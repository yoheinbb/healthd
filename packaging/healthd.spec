Name:           healthd
Version:        %{version}
Release:        %{release}%{?dist}
Summary:        Health check daemon
License:        MIT
URL:            https://github.com/yoheinbb/healthd
Source0:        healthd

BuildArch:      %{_target_cpu}

%description
healthd executes a health check script and exposes the result over HTTP.

%prep

%build

%install
install -D -m 0755 %{SOURCE0} %{buildroot}/usr/bin/healthd

%files
/usr/bin/healthd

%changelog
* Thu Apr 30 2026 healthd CI <noreply@example.com> - %{version}-%{release}
- Build rpm in CI