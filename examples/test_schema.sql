-- pgEdge Anonymizer Test Schema
-- A fictional e-commerce/CRM application with extensive PII

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Lookup tables
CREATE TABLE public.countries (
    id SERIAL PRIMARY KEY,
    code CHAR(2) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE public.states (
    id SERIAL PRIMARY KEY,
    country_id INTEGER NOT NULL REFERENCES public.countries(id),
    code VARCHAR(10) NOT NULL,
    name VARCHAR(100) NOT NULL,
    UNIQUE(country_id, code)
);

CREATE TABLE public.categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id INTEGER REFERENCES public.categories(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Core customer table with PII
CREATE TABLE public.customers (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    mobile_phone VARCHAR(50),
    date_of_birth DATE,
    ssn VARCHAR(20),
    passport_number VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    loyalty_tier VARCHAR(20) DEFAULT 'bronze'
);

CREATE INDEX idx_customers_email ON public.customers(email);
CREATE INDEX idx_customers_name ON public.customers(last_name, first_name);
CREATE INDEX idx_customers_dob ON public.customers(date_of_birth);

-- Customer addresses (one customer can have multiple)
CREATE TABLE public.addresses (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES public.customers(id) ON DELETE CASCADE,
    address_type VARCHAR(20) NOT NULL DEFAULT 'home', -- home, work, shipping, billing
    street_address TEXT NOT NULL,
    apartment VARCHAR(50),
    city VARCHAR(100) NOT NULL,
    state_id INTEGER REFERENCES public.states(id),
    postal_code VARCHAR(20),
    country_id INTEGER REFERENCES public.countries(id),
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_addresses_customer ON public.addresses(customer_id);

-- Employees (internal users with PII)
CREATE TABLE public.employees (
    id SERIAL PRIMARY KEY,
    employee_number VARCHAR(20) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    date_of_birth DATE,
    ssn VARCHAR(20),
    hire_date DATE NOT NULL,
    department VARCHAR(100),
    job_title VARCHAR(100),
    salary DECIMAL(12,2),
    manager_id INTEGER REFERENCES public.employees(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_employees_email ON public.employees(email);
CREATE INDEX idx_employees_dept ON public.employees(department);

-- Products catalog
CREATE TABLE public.products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES public.categories(id),
    price DECIMAL(10,2) NOT NULL,
    stock_quantity INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_products_category ON public.products(category_id);
CREATE INDEX idx_products_sku ON public.products(sku);

-- Orders
CREATE TABLE public.orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id INTEGER NOT NULL REFERENCES public.customers(id),
    employee_id INTEGER REFERENCES public.employees(id), -- sales rep
    shipping_address_id INTEGER REFERENCES public.addresses(id),
    billing_address_id INTEGER REFERENCES public.addresses(id),
    status VARCHAR(20) DEFAULT 'pending',
    subtotal DECIMAL(12,2),
    tax DECIMAL(10,2),
    shipping_cost DECIMAL(10,2),
    total DECIMAL(12,2),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_orders_customer ON public.orders(customer_id);
CREATE INDEX idx_orders_status ON public.orders(status);
CREATE INDEX idx_orders_date ON public.orders(created_at);

-- Order items
CREATE TABLE public.order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES public.orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES public.products(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2) DEFAULT 0,
    total DECIMAL(12,2)
);

CREATE INDEX idx_order_items_order ON public.order_items(order_id);
CREATE INDEX idx_order_items_product ON public.order_items(product_id);

-- Payment methods (credit cards with PII)
CREATE TABLE public.payment_methods (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES public.customers(id) ON DELETE CASCADE,
    card_type VARCHAR(20), -- visa, mastercard, amex
    card_number VARCHAR(20) NOT NULL, -- stored encrypted in real app
    card_expiry VARCHAR(10) NOT NULL,
    card_cvv VARCHAR(10), -- would never store this in real app!
    cardholder_name VARCHAR(200) NOT NULL,
    billing_address_id INTEGER REFERENCES public.addresses(id),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_payment_methods_customer ON public.payment_methods(customer_id);

-- Payments for orders
CREATE TABLE public.payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES public.orders(id),
    payment_method_id INTEGER REFERENCES public.payment_methods(id),
    amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    transaction_id VARCHAR(100),
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_payments_order ON public.payments(order_id);

-- Customer support tickets
CREATE TABLE public.support_tickets (
    id SERIAL PRIMARY KEY,
    ticket_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id INTEGER NOT NULL REFERENCES public.customers(id),
    assigned_to INTEGER REFERENCES public.employees(id),
    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL, -- may contain PII in free text
    status VARCHAR(20) DEFAULT 'open',
    priority VARCHAR(20) DEFAULT 'medium',
    created_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP
);

CREATE INDEX idx_tickets_customer ON public.support_tickets(customer_id);
CREATE INDEX idx_tickets_status ON public.support_tickets(status);

-- Ticket comments/notes
CREATE TABLE public.ticket_comments (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL REFERENCES public.support_tickets(id) ON DELETE CASCADE,
    author_type VARCHAR(20) NOT NULL, -- customer, employee
    author_id INTEGER NOT NULL,
    comment TEXT NOT NULL, -- may contain PII
    is_internal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_ticket_comments_ticket ON public.ticket_comments(ticket_id);

-- Customer notes (internal CRM notes)
CREATE TABLE public.customer_notes (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES public.customers(id) ON DELETE CASCADE,
    employee_id INTEGER NOT NULL REFERENCES public.employees(id),
    note TEXT NOT NULL, -- contains PII discussions
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_customer_notes_customer ON public.customer_notes(customer_id);

-- Audit log (tracks changes with user info)
CREATE TABLE public.audit_log (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id INTEGER NOT NULL,
    action VARCHAR(20) NOT NULL, -- INSERT, UPDATE, DELETE
    old_values JSONB,
    new_values JSONB,
    changed_by_type VARCHAR(20), -- customer, employee
    changed_by_id INTEGER,
    changed_by_name VARCHAR(200), -- denormalized for history
    changed_by_email VARCHAR(255),
    ip_address VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_audit_table ON public.audit_log(table_name, record_id);
CREATE INDEX idx_audit_date ON public.audit_log(created_at);
