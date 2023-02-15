#!/bin/sh

# TODO: Make these empty if we are not connected to compatible terminal
reset="\033[0m"
red="\033[31m"
green="\033[32m"
yellow="\033[1;33m"

service_name="aperture-agent"
config_dir="/etc/aperture/aperture-agent/config/"
agent_config="${config_dir}/aperture-agent.yaml"
example_config="${agent_config}.example"

has_systemctl() {
  command -V systemctl >/dev/null 2>&1
}

systemctl_version() {
  if has_systemctl; then
    systemctl --version | head -1 | sed 's/systemd //g'
  else
    printf '%b Is not a systemctl-compatible system! %b\n' "${red}" "${reset}"
    return 1
  fi
}

cleanup() {
  # This is where you remove files that were not needed on this platform / system
  if ! has_systemctl; then
    rm -f /path/to/"${service_name}".service
  fi
}

commonInstall() {
  if has_systemctl; then
    printf "%b Reload the service unit from disk%b\n" "${green}" "${reset}"
    systemctl daemon-reload ||:
  fi
  if ! [ -e "${agent_config}" ]; then
    printf "%b Installing example config - please configure the connection to etcd and prometheus %b\n" "${yellow}" "${reset}"
    cp "${example_config}" "${agent_config}"
  fi
}

cleanInstall() {
  printf "%bPost Install of an clean install%b\n" "${green}" "${reset}"
  commonInstall
  # Step 3 (clean install), enable the service in the proper way for this platform
  if has_systemctl; then
    # # rhel/centos7 cannot use ExecStartPre=+ to specify the pre start should be run as root
    # # even if you want your service to run as non root.
    # if [ "${systemd_version}" -lt 231 ]; then
    #     printf "%b systemd version %b is less then 231, fixing the service file %b\n" "${red}" "${systemd_version}" "${reset}"
    #     sed -i "s/=+/=/g" /path/to/"${service_name}".service
    # fi
    printf "%b Unmask the service%b\n" "${green}" "${reset}"
    systemctl unmask "${service_name}" ||:
    printf "%b Set the preset flag for the service unit%b\n" "${green}" "${reset}"
    systemctl preset "${service_name}" ||:
    # We don't want to enable the service by default,
    # as we require some changes to config files to be made
    # printf "%b Set the enabled flag for the service unit%b\n" "${green}" "${reset}"
    # systemctl enable "${service_name}" ||:
    # systemctl restart "${service_name}" ||:
  fi
}

upgrade() {
  printf "%b Post Install of an upgrade%b\n" "${green}" "${reset}"
  commonInstall
}

# Step 2, check if this is a clean install or an upgrade
action="$1"
if  [ "$1" = "configure" ] && [ -z "$2" ]; then
  # Alpine linux does not pass args, and deb passes $1=configure
  action="install"
elif [ "$1" = "configure" ] && [ -n "$2" ]; then
  # deb passes $1=configure $2=<current version>
  action="upgrade"
fi

case "$action" in
  "1" | "install")
    cleanInstall
    ;;
  "2" | "upgrade")
    upgrade
    ;;
  *)
    # $1 == version being installed
    cleanInstall
    ;;
esac

# Step 4, clean up unused files, yes you get a warning when you remove the package, but that is ok.
cleanup
