# /!\ This file is maintained at https://github.com/openshift/node_exporter
%global debug_package   %{nil}

%global provider        github
%global provider_tld    com
%global project         prometheus
%global repo            node_exporter
# https://github.com/prometheus/node_exporter
%global provider_prefix %{provider}.%{provider_tld}/%{project}/%{repo}
%global import_path     %{provider_prefix}
# %commit is intended to be set by tito. The values in this spec file will not be kept up to date.
%{!?commit:
%global commit          774a8c9ce7df6ce20b12be3b3f2175a74e8ddb21
}
%global shortcommit     %(c=%{commit}; echo ${c:0:7})
%global gopathdir       %{_sourcedir}/go
%global upstream_ver    0.17.0
%global rpm_ver         %(v=%{upstream_ver}; echo ${v//-/_})
%global download_prefix %{provider}.%{provider_tld}/openshift/%{repo}

Name:		golang-%{provider}-%{project}-%{repo}
Version:	%{rpm_ver}
Release:	1.git%{shortcommit}%{?dist}
Summary:	Prometheus exporter for hardware and OS metrics exposed by *NIX kernels
License:	ASL 2.0
URL:		https://prometheus.io/
Source0:	https://%{download_prefix}/archive/%{commit}/%{repo}-%{commit}.tar.gz

# e.g. el6 has ppc64 arch without gcc-go, so EA tag is required
ExclusiveArch:  %{?go_arches:%{go_arches}}%{!?go_arches:%{ix86} x86_64 aarch64 %{arm} ppc64le s390x}
# If go_compiler is not set to 1, there is no virtual provide. Use golang instead.
BuildRequires: %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}
BuildRequires: glibc-static
BuildRequires: prometheus-promu

%description
%{summary}

%package -n %{project}-node-exporter
Summary:        %{summary}
Provides:       prometheus-node-exporter = %{version}-%{release}
Provides:       prometheus-node_exporter = %{version}-%{release}
Obsoletes:      prometheus-node_exporter < %{version}-%{release}

%description -n %{project}-node-exporter
%{summary}

%prep
%setup -q -n %{repo}-%{commit}

%build
# Go expects a full path to the sources which is not included in the source
# tarball so create a link with the expected path
mkdir -p %{gopathdir}/src/%{provider}.%{provider_tld}/%{project}
GOSRCDIR=%{gopathdir}/src/%{import_path}
if [ ! -e "$GOSRCDIR" ]; then
  ln -s `pwd` "$GOSRCDIR"
fi
export GOPATH=%{gopathdir}

make build BUILD_PROMU=false

%install
install -d %{buildroot}%{_bindir}
install -D -p -m 0755 node_exporter %{buildroot}%{_bindir}/node_exporter
ln -s %{_bindir}/node_exporter \
      %{buildroot}%{_bindir}/prometheus-node-exporter
install -D -p -m 0644 prometheus-node-exporter.service \
                      %{buildroot}%{_unitdir}/prometheus-node-exporter.service
install -D -p -m 0644 prometheus-node-exporter.sysconfig \
                      %{buildroot}%{_sysconfdir}/sysconfig/prometheus-node-exporter

%files -n %{project}-node-exporter
%license LICENSE NOTICE
%doc CHANGELOG.md CONTRIBUTING.md MAINTAINERS.md README.md
%{_bindir}/node_exporter
%{_bindir}/prometheus-node-exporter
%{_unitdir}/prometheus-node-exporter.service
%{_sysconfdir}/sysconfig/prometheus-node-exporter

%changelog
* Tue Jan 15 2019 Paul Gier <pgier@redhat.com> - 0.17.0-1
- upgrade to 0.17.0

* Thu Sep 27 2018 Simon Pasquier <spasqui@redhat.com> - 0.16.0-3
- Fix stop command in systemd unit

* Fri Jul 27 2018 Simon Pasquier <spasquie@redhat.com> - 0.16.0-2
- Enable aarch64

* Wed Jul 25 2018 Paul Gier <pgier@redhat.com> - 0.16.0-1
- upgrade to 0.16.0

* Wed Jun 20 2018 Simon Pasquier <spasquie@redhat.com> - 0.15.2-3
- Fix systemd configuration

* Tue May 22 2018 Paul Gier <pgier@redhat.com> - 0.15.2-2
- Add systemd unit file and related config

* Wed Jan 17 2018 Paul Gier <pgier@redhat.com> - 0.15.2-1
- upgrade to 0.15.2

* Wed Jan 10 2018 Yaakov Selkowitz <yselkowi@redhat.com> - 0.15.1-2
- Rebuilt for ppc64le, s390x enablement

* Thu Nov 09 2017 Paul Gier <pgier@redhat.com> - 0.15.1-1
- upgrade to 0.15.1
- added build requires on glibc-static

* Wed Aug 30 2017 Paul Gier <pgier@redhat.com> - 0.14.0-1
- Initial package creation


