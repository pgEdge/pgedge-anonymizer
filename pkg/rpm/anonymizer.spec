%global debug_package %{nil}
%global sname pgedge-anonymizer

Name:           %{sname}
Version:        %{anonymizer_version}
Release:        %{anonymizer_buildnum}%{?dist}
Summary:        pgEdge Anonymizer is a tool for anonymizing Personally Identifiable Information (PII) in PostgreSQL databases

License:        PostgreSQL License
URL:            https://github.com/pgEdge/%{sname}

Source0:	https://github.com/pgEdge/%{sname}/releases/download/v%{anonymizer_version}/%{sname}_%{anonymizer_version}_Linux_%{arch}.tar.gz
Source1:        %{sname}.yaml
Source2:	LICENCE.md
Source3:        %{sname}-patterns.yaml

%description
pgEdge Anonymizer is a tool for anonymizing Personally Identifiable
Information (PII) in PostgreSQL databases to meet GDPR and other regulatory
requirements when cloning production data for development purposes.

The tool processes columns specified in a configuration file, applying
pattern-based anonymization while maintaining data consistency across tables

%prep
%setup -q -c -n pgedge-anonymizer-%{version}
cp %{SOURCE2} .

%build
syft dir:%{_builddir} -o cyclonedx-json > %{_builddir}/%{sname}-sbom.json || exit 1

KEY_ID=$(gpg --list-secret-keys --with-colons | awk -F: '/^sec/{print $5}' | head -n 1); export KEY_ID
gpg --armor --detach-sign --local-user "$KEY_ID" --output %{_builddir}/%{sname}-sbom.json.asc %{_builddir}/%{sname}-sbom.json || exit 1

%install
install -D -m 0755 %{sname} %{buildroot}/usr/bin/%{sname}
mkdir -p %{buildroot}%{_sysconfdir}/pgedge
install -D -m 0644 %{SOURCE1} %{buildroot}%{_sysconfdir}/pgedge/%{sname}.yaml
install -D -m 0644 %{SOURCE3} %{buildroot}%{_sysconfdir}/pgedge/%{sname}-patterns.yaml
mkdir -p %{buildroot}%{_datadir}/%{sname}
install -p -m 0644 %{_builddir}/%{sname}-sbom.json %{buildroot}%{_datadir}/%{sname}/%{sname}-sbom.json
install -p -m 0644 %{_builddir}/%{sname}-sbom.json.asc %{buildroot}%{_datadir}/%{sname}/%{sname}-sbom.json.asc

%files
%license LICENCE.md
%doc README.md
%{_bindir}/%{sname}
%dir %{_sysconfdir}/pgedge
%config(noreplace) %{_sysconfdir}/pgedge/%{sname}.yaml
%config(noreplace) %{_sysconfdir}/pgedge/%{sname}-patterns.yaml
%{_datadir}/%{sname}/%{sname}-sbom.json
%{_datadir}/%{sname}/%{sname}-sbom.json.asc

%changelog
* Thu Apr 02 2026 Muhammad Aqeel <muhammad.aqeel@pgedge.com> - 1.0.0
- Update RPM package of pgedge-anonymizer
* Mon Dec 22 2025 Muhammad Aqeel <muhammad.aqeel@pgedge.com> - 1.0.0-beta2
- Update RPM package of pgedge-anonymizer
* Mon Dec 15 2025 Muhammad Aqeel <muhammad.aqeel@pgedge.com> - 1.0.0-beta1
- Initial RPM package of pgedge-anonymizer

