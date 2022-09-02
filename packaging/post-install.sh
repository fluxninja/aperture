#!/bin/sh

# TODO: Make these empty if we are not connected to compatible terminal
reset="\033[0m"
red="\033[31m"
green="\033[32m"

service_name="aperture-agent"

has_systemctl() {
  command -V systemctl >/dev/null 2>&1
}

systemctl_version() {
  if has_systemctl; then
    systemctl --version | head -1 | sed 's/systemd //g'
  else
    printf '%s Is not a systemctl-compatible system! %s\n' "${red}" "${reset}"
    return 1
  fi
}

cleanup() {
  # This is where you remove files that were not needed on this platform / system
  if ! has_systemctl; then
    rm -f /path/to/"${service_name}".service
  fi
}

cleanInstall() {
  printf "%sPost Install of an clean install%s\n" "${green}" "${reset}"
  # Step 3 (clean install), enable the service in the proper way for this platform
  if has_systemctl; then
    # # rhel/centos7 cannot use ExecStartPre=+ to specify the pre start should be run as root
    # # even if you want your service to run as non root.
    # if [ "${systemd_version}" -lt 231 ]; then
    #     printf "%s systemd version %s is less then 231, fixing the service file %s\n" "${red}" "${systemd_version}" "${reset}"
    #     sed -i "s/=+/=/g" /path/to/"${service_name}".service
    # fi
    printf "%s Reload the service unit from disk%s\n" "${green}" "${reset}"
    systemctl daemon-reload ||:
    printf "%s Unmask the service%s\n" "${green}" "${reset}"
    systemctl unmask "${service_name}" ||:
    printf "%s Set the preset flag for the service unit%s\n" "${green}" "${reset}"
    systemctl preset "${service_name}" ||:
    # We don't want to enable the service by default,
    # as we require some changes to config files to be made
    # printf "%s Set the enabled flag for the service unit%s\n" "${green}" "${reset}"
    # systemctl enable "${service_name}" ||:
    # systemctl restart "${service_name}" ||:
  fi
}

upgrade() {
    printf "%s Post Install of an upgrade%s\n" "${green}" "${reset}"
    # Step 3(upgrade), do what you need
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
