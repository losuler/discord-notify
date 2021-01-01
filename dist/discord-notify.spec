Name:           discord-notify
Version:        0.1.0
Release:        1%{?dist}
Summary:        Receive notifications of new messages from Discord through Telegram.
License:        GPLv3
URL:            https://gitlab.com/losuler/%{name}
Source0:        https://gitlab.com/losuler/%{name}/-/archive/v%{version}/%{name}-v%{version}.tar.gz

BuildRequires:  golang >= 1.12
BuildRequires:  systemd

%description
Provides a simple way to recieve notifications of new messages from Discord through Telegram. 
Once a message is recieved on Discord the username of the sender is sent to Telegram through 
the Telegram Bot API.

%prep
%setup -q -n %{name}-v%{version}

%build
go build \
    -ldflags "${LDFLAGS:-} -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')%{?__global_ldflags: -extldflags '%__global_ldflags'}" \
    -a -v

%install
install -D -m 0755 %{name} %{buildroot}/%{_bindir}/%{name}
install -D -m 0640 conf.json.example %{buildroot}/%{_sysconfdir}/%{name}.json
install -D -m 0644 dist/%{name}.service %{buildroot}/%{_unitdir}/%{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%files
%license LICENSE.txt
%doc README.md
%{_bindir}/%{name}
%config %{_sysconfdir}/%{name}.json
%{_unitdir}/%{name}.service

%changelog
