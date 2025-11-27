-- pgEdge Anonymizer Test Data Population
-- Generates realistic dummy PII data

-- Helper function for random integers
CREATE OR REPLACE FUNCTION random_between(low INT, high INT)
RETURNS INT AS $$
BEGIN
    RETURN floor(random() * (high - low + 1) + low);
END;
$$ LANGUAGE plpgsql;

-- Helper function for random element from array
CREATE OR REPLACE FUNCTION random_element(arr TEXT[])
RETURNS TEXT AS $$
BEGIN
    RETURN arr[1 + floor(random() * array_length(arr, 1))::int];
END;
$$ LANGUAGE plpgsql;

-- Populate countries
INSERT INTO public.countries (code, name) VALUES
('US', 'United States'),
('CA', 'Canada'),
('GB', 'United Kingdom'),
('AU', 'Australia'),
('DE', 'Germany');

-- Populate US states
INSERT INTO public.states (country_id, code, name) VALUES
(1, 'CA', 'California'),
(1, 'NY', 'New York'),
(1, 'TX', 'Texas'),
(1, 'FL', 'Florida'),
(1, 'IL', 'Illinois'),
(1, 'PA', 'Pennsylvania'),
(1, 'OH', 'Ohio'),
(1, 'GA', 'Georgia'),
(1, 'NC', 'North Carolina'),
(1, 'MI', 'Michigan'),
(1, 'WA', 'Washington'),
(1, 'AZ', 'Arizona'),
(1, 'MA', 'Massachusetts'),
(1, 'CO', 'Colorado'),
(1, 'VA', 'Virginia');

-- Populate categories
INSERT INTO public.categories (name, description, parent_id) VALUES
('Electronics', 'Electronic devices and accessories', NULL),
('Clothing', 'Apparel and fashion', NULL),
('Home & Garden', 'Home improvement and gardening', NULL),
('Sports', 'Sports equipment and gear', NULL),
('Books', 'Books and publications', NULL);

INSERT INTO public.categories (name, description, parent_id) VALUES
('Smartphones', 'Mobile phones', 1),
('Laptops', 'Portable computers', 1),
('Cameras', 'Digital cameras', 1),
('Men''s Clothing', 'Men''s apparel', 2),
('Women''s Clothing', 'Women''s apparel', 2),
('Furniture', 'Home furniture', 3),
('Garden Tools', 'Gardening equipment', 3);

-- Arrays for generating realistic data
DO $$
DECLARE
    first_names TEXT[] := ARRAY[
        'James', 'Mary', 'John', 'Patricia', 'Robert', 'Jennifer', 'Michael', 'Linda',
        'William', 'Elizabeth', 'David', 'Barbara', 'Richard', 'Susan', 'Joseph', 'Jessica',
        'Thomas', 'Sarah', 'Charles', 'Karen', 'Christopher', 'Nancy', 'Daniel', 'Lisa',
        'Matthew', 'Betty', 'Anthony', 'Margaret', 'Mark', 'Sandra', 'Donald', 'Ashley',
        'Steven', 'Kimberly', 'Paul', 'Emily', 'Andrew', 'Donna', 'Joshua', 'Michelle',
        'Kenneth', 'Dorothy', 'Kevin', 'Carol', 'Brian', 'Amanda', 'George', 'Melissa',
        'Edward', 'Deborah', 'Ronald', 'Stephanie', 'Timothy', 'Rebecca', 'Jason', 'Sharon',
        'Jeffrey', 'Laura', 'Ryan', 'Cynthia', 'Jacob', 'Kathleen', 'Gary', 'Amy',
        'Nicholas', 'Angela', 'Eric', 'Shirley', 'Jonathan', 'Anna', 'Stephen', 'Brenda',
        'Larry', 'Pamela', 'Justin', 'Emma', 'Scott', 'Nicole', 'Brandon', 'Helen',
        'Benjamin', 'Samantha', 'Samuel', 'Katherine', 'Raymond', 'Christine', 'Gregory', 'Debra',
        'Frank', 'Rachel', 'Alexander', 'Carolyn', 'Patrick', 'Janet', 'Jack', 'Catherine'
    ];

    last_names TEXT[] := ARRAY[
        'Smith', 'Johnson', 'Williams', 'Brown', 'Jones', 'Garcia', 'Miller', 'Davis',
        'Rodriguez', 'Martinez', 'Hernandez', 'Lopez', 'Gonzalez', 'Wilson', 'Anderson', 'Thomas',
        'Taylor', 'Moore', 'Jackson', 'Martin', 'Lee', 'Perez', 'Thompson', 'White',
        'Harris', 'Sanchez', 'Clark', 'Ramirez', 'Lewis', 'Robinson', 'Walker', 'Young',
        'Allen', 'King', 'Wright', 'Scott', 'Torres', 'Nguyen', 'Hill', 'Flores',
        'Green', 'Adams', 'Nelson', 'Baker', 'Hall', 'Rivera', 'Campbell', 'Mitchell',
        'Carter', 'Roberts', 'Turner', 'Phillips', 'Evans', 'Parker', 'Edwards', 'Collins',
        'Stewart', 'Morris', 'Murphy', 'Cook', 'Rogers', 'Morgan', 'Peterson', 'Cooper',
        'Reed', 'Bailey', 'Bell', 'Gomez', 'Kelly', 'Howard', 'Ward', 'Cox',
        'Diaz', 'Richardson', 'Wood', 'Watson', 'Brooks', 'Bennett', 'Gray', 'James',
        'Reyes', 'Cruz', 'Hughes', 'Price', 'Myers', 'Long', 'Foster', 'Sanders',
        'Ross', 'Morales', 'Powell', 'Sullivan', 'Russell', 'Ortiz', 'Jenkins', 'Gutierrez'
    ];

    street_names TEXT[] := ARRAY[
        'Main', 'Oak', 'Maple', 'Cedar', 'Pine', 'Elm', 'Washington', 'Lake',
        'Hill', 'Park', 'River', 'Spring', 'Valley', 'Forest', 'Sunset', 'Highland',
        'Meadow', 'Church', 'Mill', 'School', 'North', 'South', 'East', 'West',
        'Brookside', 'Lakeview', 'Hillcrest', 'Fairview', 'Woodland', 'Greenwood', 'Lincoln', 'Franklin'
    ];

    street_types TEXT[] := ARRAY[
        'Street', 'Avenue', 'Boulevard', 'Drive', 'Lane', 'Road', 'Way', 'Court',
        'Place', 'Circle', 'Terrace', 'Trail'
    ];

    cities TEXT[] := ARRAY[
        'Los Angeles', 'New York', 'Chicago', 'Houston', 'Phoenix', 'Philadelphia',
        'San Antonio', 'San Diego', 'Dallas', 'San Jose', 'Austin', 'Jacksonville',
        'Fort Worth', 'Columbus', 'Charlotte', 'San Francisco', 'Indianapolis', 'Seattle',
        'Denver', 'Boston', 'Nashville', 'Portland', 'Las Vegas', 'Detroit', 'Memphis',
        'Louisville', 'Milwaukee', 'Albuquerque', 'Tucson', 'Fresno', 'Sacramento', 'Atlanta'
    ];

    email_domains TEXT[] := ARRAY[
        'gmail.com', 'yahoo.com', 'hotmail.com', 'outlook.com', 'icloud.com',
        'aol.com', 'mail.com', 'protonmail.com', 'fastmail.com', 'zoho.com'
    ];

    departments TEXT[] := ARRAY[
        'Sales', 'Marketing', 'Engineering', 'Customer Support', 'Human Resources',
        'Finance', 'Operations', 'Product', 'Legal', 'IT'
    ];

    job_titles TEXT[] := ARRAY[
        'Manager', 'Senior Specialist', 'Specialist', 'Associate', 'Director',
        'Analyst', 'Coordinator', 'Representative', 'Lead', 'Consultant'
    ];

    product_adjectives TEXT[] := ARRAY[
        'Premium', 'Professional', 'Essential', 'Classic', 'Advanced', 'Ultra',
        'Elite', 'Basic', 'Deluxe', 'Standard'
    ];

    product_nouns TEXT[] := ARRAY[
        'Widget', 'Gadget', 'Device', 'Tool', 'System', 'Kit', 'Set', 'Pack',
        'Bundle', 'Collection', 'Series', 'Edition'
    ];

    ticket_subjects TEXT[] := ARRAY[
        'Order not received', 'Wrong item delivered', 'Refund request', 'Product quality issue',
        'Account access problem', 'Billing question', 'Shipping delay', 'Return request',
        'Product information needed', 'Technical support', 'Password reset', 'Payment failed'
    ];

    i INT;
    j INT;
    fname TEXT;
    lname TEXT;
    cust_id INT;
    emp_id INT;
    order_id INT;
    ticket_id INT;
    addr_id INT;
    payment_method_id INT;
    product_count INT;
    v_phone TEXT;
    v_ssn TEXT;
    v_card TEXT;
    v_dob DATE;
    v_state_id INT;

BEGIN
    -- Generate 500 customers
    FOR i IN 1..500 LOOP
        fname := first_names[1 + floor(random() * array_length(first_names, 1))::int];
        lname := last_names[1 + floor(random() * array_length(last_names, 1))::int];
        v_state_id := 1 + floor(random() * 15)::int;
        v_dob := '1950-01-01'::date + (random() * 20000)::int;
        v_phone := '(' || (200 + floor(random() * 800)::int)::text || ') ' ||
                   (200 + floor(random() * 800)::int)::text || '-' ||
                   lpad((floor(random() * 10000)::int)::text, 4, '0');
        v_ssn := lpad((floor(random() * 900)::int + 100)::text, 3, '0') || '-' ||
                 lpad((floor(random() * 100)::int)::text, 2, '0') || '-' ||
                 lpad((floor(random() * 10000)::int)::text, 4, '0');

        INSERT INTO public.customers (
            first_name, last_name, email, phone, mobile_phone, date_of_birth, ssn,
            passport_number, loyalty_tier, created_at
        ) VALUES (
            fname,
            lname,
            lower(fname || '.' || lname || i::text || '@' ||
                  email_domains[1 + floor(random() * array_length(email_domains, 1))::int]),
            v_phone,
            CASE WHEN random() > 0.3 THEN
                '(' || (200 + floor(random() * 800)::int)::text || ') ' ||
                (200 + floor(random() * 800)::int)::text || '-' ||
                lpad((floor(random() * 10000)::int)::text, 4, '0')
            ELSE NULL END,
            v_dob,
            CASE WHEN random() > 0.7 THEN v_ssn ELSE NULL END,
            CASE WHEN random() > 0.9 THEN
                chr(65 + floor(random() * 26)::int) || lpad((floor(random() * 100000000)::int)::text, 8, '0')
            ELSE NULL END,
            (ARRAY['bronze', 'silver', 'gold', 'platinum'])[1 + floor(random() * 4)::int],
            NOW() - (random() * 1000)::int * interval '1 day'
        ) RETURNING id INTO cust_id;

        -- Add 1-3 addresses per customer
        FOR j IN 1..random_between(1, 3) LOOP
            INSERT INTO public.addresses (
                customer_id, address_type, street_address, apartment, city,
                state_id, postal_code, country_id, is_primary
            ) VALUES (
                cust_id,
                (ARRAY['home', 'work', 'shipping', 'billing'])[j],
                (100 + floor(random() * 9900)::int)::text || ' ' ||
                    street_names[1 + floor(random() * array_length(street_names, 1))::int] || ' ' ||
                    street_types[1 + floor(random() * array_length(street_types, 1))::int],
                CASE WHEN random() > 0.7 THEN 'Apt ' || (100 + floor(random() * 900)::int)::text ELSE NULL END,
                cities[1 + floor(random() * array_length(cities, 1))::int],
                v_state_id,
                lpad((10000 + floor(random() * 90000)::int)::text, 5, '0'),
                1,
                j = 1
            );
        END LOOP;
    END LOOP;

    -- Generate 50 employees
    FOR i IN 1..50 LOOP
        fname := first_names[1 + floor(random() * array_length(first_names, 1))::int];
        lname := last_names[1 + floor(random() * array_length(last_names, 1))::int];
        v_dob := '1960-01-01'::date + (random() * 15000)::int;
        v_ssn := lpad((floor(random() * 900)::int + 100)::text, 3, '0') || '-' ||
                 lpad((floor(random() * 100)::int)::text, 2, '0') || '-' ||
                 lpad((floor(random() * 10000)::int)::text, 4, '0');

        INSERT INTO public.employees (
            employee_number, first_name, last_name, email, phone,
            date_of_birth, ssn, hire_date, department, job_title, salary,
            manager_id
        ) VALUES (
            'EMP' || lpad(i::text, 5, '0'),
            fname,
            lname,
            lower(fname || '.' || lname || '@company.com'),
            '(' || (200 + floor(random() * 800)::int)::text || ') ' ||
                (200 + floor(random() * 800)::int)::text || '-' ||
                lpad((floor(random() * 10000)::int)::text, 4, '0'),
            v_dob,
            v_ssn,
            '2015-01-01'::date + (random() * 3000)::int,
            departments[1 + floor(random() * array_length(departments, 1))::int],
            job_titles[1 + floor(random() * array_length(job_titles, 1))::int],
            40000 + floor(random() * 100000)::int,
            CASE WHEN i > 5 THEN 1 + floor(random() * 5)::int ELSE NULL END
        );
    END LOOP;

    -- Generate 200 products
    FOR i IN 1..200 LOOP
        INSERT INTO public.products (
            sku, name, description, category_id, price, stock_quantity
        ) VALUES (
            'SKU' || lpad(i::text, 6, '0'),
            product_adjectives[1 + floor(random() * array_length(product_adjectives, 1))::int] || ' ' ||
                product_nouns[1 + floor(random() * array_length(product_nouns, 1))::int] || ' ' || i::text,
            'High quality product with excellent features and specifications. Perfect for everyday use.',
            6 + floor(random() * 7)::int,
            9.99 + floor(random() * 500)::numeric,
            floor(random() * 1000)::int
        );
    END LOOP;

    -- Generate payment methods for customers (about 60%)
    FOR cust_id IN SELECT id FROM public.customers WHERE random() > 0.4 LOOP
        -- Get a billing address for this customer
        SELECT id INTO addr_id FROM public.addresses WHERE customer_id = cust_id LIMIT 1;

        v_card := '4' || lpad((floor(random() * 1000000000000000)::bigint)::text, 15, '0');

        INSERT INTO public.payment_methods (
            customer_id, card_type, card_number, card_expiry, card_cvv,
            cardholder_name, billing_address_id, is_default
        )
        SELECT
            cust_id,
            (ARRAY['visa', 'mastercard', 'amex', 'discover'])[1 + floor(random() * 4)::int],
            v_card,
            lpad((1 + floor(random() * 12)::int)::text, 2, '0') || '/' ||
                (25 + floor(random() * 6)::int)::text,
            lpad((floor(random() * 1000)::int)::text, 3, '0'),
            c.first_name || ' ' || c.last_name,
            addr_id,
            TRUE
        FROM public.customers c WHERE c.id = cust_id;
    END LOOP;

    -- Generate 2000 orders
    FOR i IN 1..2000 LOOP
        -- Select random customer
        SELECT id INTO cust_id FROM public.customers ORDER BY random() LIMIT 1;
        SELECT id INTO emp_id FROM public.employees ORDER BY random() LIMIT 1;
        SELECT id INTO addr_id FROM public.addresses WHERE customer_id = cust_id LIMIT 1;

        INSERT INTO public.orders (
            order_number, customer_id, employee_id, shipping_address_id,
            billing_address_id, status, subtotal, tax, shipping_cost, total,
            notes, created_at
        ) VALUES (
            'ORD' || lpad(i::text, 8, '0'),
            cust_id,
            CASE WHEN random() > 0.5 THEN emp_id ELSE NULL END,
            addr_id,
            addr_id,
            (ARRAY['pending', 'processing', 'shipped', 'delivered', 'cancelled'])[
                1 + floor(random() * 5)::int],
            0, -- will update
            0,
            5 + floor(random() * 20)::numeric,
            0,
            CASE WHEN random() > 0.8 THEN 'Customer requested gift wrapping' ELSE NULL END,
            NOW() - (random() * 365)::int * interval '1 day'
        ) RETURNING id INTO order_id;

        -- Add 1-5 items per order
        product_count := 1 + floor(random() * 5)::int;
        INSERT INTO public.order_items (order_id, product_id, quantity, unit_price, total)
        SELECT
            order_id,
            p.id,
            1 + floor(random() * 3)::int as qty,
            p.price,
            (1 + floor(random() * 3)::int) * p.price
        FROM public.products p
        ORDER BY random()
        LIMIT product_count;

        -- Update order totals
        UPDATE public.orders SET
            subtotal = (SELECT COALESCE(SUM(total), 0) FROM public.order_items WHERE order_items.order_id = orders.id),
            tax = (SELECT COALESCE(SUM(total), 0) * 0.08 FROM public.order_items WHERE order_items.order_id = orders.id),
            total = (SELECT COALESCE(SUM(total), 0) * 1.08 FROM public.order_items WHERE order_items.order_id = orders.id) + shipping_cost
        WHERE id = order_id;

        -- Add payment for non-cancelled orders
        IF random() > 0.2 THEN
            SELECT id INTO payment_method_id FROM public.payment_methods WHERE customer_id = cust_id LIMIT 1;
            INSERT INTO public.payments (order_id, payment_method_id, amount, status, transaction_id, processed_at)
            SELECT order_id, payment_method_id, total,
                   (ARRAY['completed', 'pending', 'failed'])[1 + floor(random() * 3)::int],
                   'TXN' || lpad((floor(random() * 1000000000)::int)::text, 12, '0'),
                   CASE WHEN random() > 0.3 THEN NOW() - (random() * 30)::int * interval '1 day' ELSE NULL END
            FROM public.orders WHERE id = order_id;
        END IF;
    END LOOP;

    -- Generate 300 support tickets
    FOR i IN 1..300 LOOP
        SELECT id INTO cust_id FROM public.customers ORDER BY random() LIMIT 1;
        SELECT id INTO emp_id FROM public.employees WHERE department = 'Customer Support' ORDER BY random() LIMIT 1;

        INSERT INTO public.support_tickets (
            ticket_number, customer_id, assigned_to, subject, description,
            status, priority, created_at
        )
        SELECT
            'TKT' || lpad(i::text, 8, '0'),
            cust_id,
            CASE WHEN random() > 0.3 THEN emp_id ELSE NULL END,
            ticket_subjects[1 + floor(random() * array_length(ticket_subjects, 1))::int],
            'Customer ' || c.first_name || ' ' || c.last_name ||
                ' (email: ' || c.email || ', phone: ' || COALESCE(c.phone, 'N/A') ||
                ') reported an issue. ' ||
                'The customer mentioned their address at ' ||
                (SELECT street_address || ', ' || city FROM public.addresses WHERE customer_id = cust_id LIMIT 1) ||
                '. Please follow up as soon as possible.',
            (ARRAY['open', 'in_progress', 'resolved', 'closed'])[1 + floor(random() * 4)::int],
            (ARRAY['low', 'medium', 'high', 'urgent'])[1 + floor(random() * 4)::int],
            NOW() - (random() * 180)::int * interval '1 day'
        FROM public.customers c WHERE c.id = cust_id
        RETURNING id INTO ticket_id;

        -- Add 1-5 comments per ticket
        FOR j IN 1..random_between(1, 5) LOOP
            INSERT INTO public.ticket_comments (ticket_id, author_type, author_id, comment, is_internal)
            VALUES (
                ticket_id,
                CASE WHEN random() > 0.5 THEN 'employee' ELSE 'customer' END,
                CASE WHEN random() > 0.5 THEN emp_id ELSE cust_id END,
                CASE
                    WHEN j = 1 THEN 'Initial contact from customer regarding their issue.'
                    WHEN random() > 0.5 THEN 'Following up on the previous communication. Customer can be reached at their phone number on file.'
                    ELSE 'Thank you for your patience. We are working on resolving this matter.'
                END,
                random() > 0.7
            );
        END LOOP;
    END LOOP;

    -- Generate customer notes
    FOR i IN 1..500 LOOP
        SELECT id INTO cust_id FROM public.customers ORDER BY random() LIMIT 1;
        SELECT id INTO emp_id FROM public.employees ORDER BY random() LIMIT 1;

        INSERT INTO public.customer_notes (customer_id, employee_id, note)
        SELECT
            cust_id,
            emp_id,
            'Spoke with ' || c.first_name || ' ' || c.last_name || ' on ' ||
            to_char(NOW() - (random() * 90)::int * interval '1 day', 'YYYY-MM-DD') ||
            '. Customer mentioned they live in ' ||
            (SELECT city FROM public.addresses WHERE customer_id = cust_id LIMIT 1) ||
            ' and prefers to be contacted at ' || COALESCE(c.phone, c.email) ||
            '. ' ||
            CASE floor(random() * 5)::int
                WHEN 0 THEN 'Very satisfied customer, potential for upsell.'
                WHEN 1 THEN 'Had some concerns about delivery times.'
                WHEN 2 THEN 'Interested in our loyalty program benefits.'
                WHEN 3 THEN 'Requested information about bulk orders.'
                ELSE 'Follow up needed regarding recent purchase.'
            END
        FROM public.customers c WHERE c.id = cust_id;
    END LOOP;

    -- Generate audit log entries
    FOR i IN 1..1000 LOOP
        SELECT id, first_name, last_name, email INTO cust_id, fname, lname, v_phone
        FROM public.customers ORDER BY random() LIMIT 1;
        SELECT id INTO emp_id FROM public.employees ORDER BY random() LIMIT 1;

        INSERT INTO public.audit_log (
            table_name, record_id, action, old_values, new_values,
            changed_by_type, changed_by_id, changed_by_name, changed_by_email, ip_address
        )
        SELECT
            (ARRAY['customers', 'orders', 'addresses', 'payments'])[1 + floor(random() * 4)::int],
            floor(random() * 500)::int + 1,
            (ARRAY['INSERT', 'UPDATE', 'DELETE'])[1 + floor(random() * 3)::int],
            '{"field": "old_value"}'::jsonb,
            '{"field": "new_value"}'::jsonb,
            CASE WHEN random() > 0.3 THEN 'employee' ELSE 'customer' END,
            CASE WHEN random() > 0.3 THEN emp_id ELSE cust_id END,
            e.first_name || ' ' || e.last_name,
            e.email,
            (floor(random() * 256)::int)::text || '.' ||
            (floor(random() * 256)::int)::text || '.' ||
            (floor(random() * 256)::int)::text || '.' ||
            (floor(random() * 256)::int)::text
        FROM public.employees e WHERE e.id = emp_id;
    END LOOP;

END $$;

-- Clean up helper functions
DROP FUNCTION IF EXISTS random_between(INT, INT);
DROP FUNCTION IF EXISTS random_element(TEXT[]);

-- Analyze tables for query optimization
ANALYZE;
