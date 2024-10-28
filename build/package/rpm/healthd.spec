Name:           healthd
Version:        0.4
Release:        1%{?dist}
Summary:        rest api server with cmd execute

License:        GPLv3
%undefine _disable_source_fetch
#Source0:        https://github.com/yoheinbb/%{name}/archive/refs/tags/v%{version}.tar.gz
Source0:        https://github.com/yoheinbb/healthd/archive/refs/tags/test.tar.gz

#BuildRequires:  golang

Provides:       %{name} = %{version}

%description
rest api server with cmd execute

%prep
%autosetup

%make_build

%install
install -Dpm 0755 build/%{name} %{buildroot}%{_prefix}/local/bin/%{name}
install -Dpm 0644 configs/conf/global.json %{buildroot}%{_sysconfdir}/%{name}/global.json
install -Dpm 0644 configs/conf/script.json %{buildroot}%{_sysconfdir}/%{name}/script.json
install -Dpm 0755 configs/scripts/sample_script %{buildroot}%{_sysconfdir}/%{name}/scripts/sample_script

%files
%{_prefix}/local/bin/%{name}
%config %{_sysconfdir}/%{name}

%changelog
* Wed Oct 23 2024 yoheinbb v0.4
- add healthd specfile

