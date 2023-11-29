# usage> Gitpod [before | init | command]
# short> run the Gitpod scripts
# long>  run the Gitpod initialisation scripts from clog
#                             _
#   ___   _ __   ___   _ _   | |_   _ __   __ _
#  / _ \ | '_ \ / -_) | ' \  |  _| | '_ \ / _` |
#  \___/ | .__/ \___| |_||_|  \__| | .__/ \__, |
#        |_|                       |_|    |___/

[ -n "$(clog -v)" ] && source <(clog Cat clogrc/core.inc.sh) ## use clog tools
case $1 in
  "before")   ACTION="before" ;;
  "init")     ACTION="init" ;;
  "command")  ACTION="command" ;;
  *)          echo -e "$cW Warning$cT specify$cC before$cT |$cC init$cT | command. Using default $cC before$cX"
              ACTION="before" ;;
esac

# ------------------------------------------------------------------------------
#   ___   ___   ___    ___    ___   ___
#  | _ ) | __| | __|  / _ \  | _ \ | __|
#  | _ \ | _|  | _|  | (_) | |   / | _|
#  |___/ |___| |_|    \___/  |_|_\ |___|
# ------------------------------------------------------------------------------
if [[ $ACTION == "before" ]] ; then

  echo -e "${Cgreen}gitpod$cC BEFORE$cT tasks"

  fnInfo "$cC yarn$cT global path"
  PATH="$PATH:$(yarn global bin)"
  echo "PATH=\"$PATH\"" >> ~/.bashrc

#   fnInfo "$cC  zmp$cT update"
#   sudo curl https://mrmxf.com/get/zmp  -o /usr/local/bin/zmp  ; sudo chmod +x /usr/local/bin/zmp

  fnInfo "$cC clog$cT update"
  curl --no-progress-meter https://mrmxf.com/get/clogbin | bash

  fnInfo "$cC  aws$cT cli install"
  curl --no-progress-meter "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "/tmp/awscliv2.zip"
  unzip -q -o /tmp/awscliv2.zip -d /tmp
  sudo /tmp/aws/install

#   fnInfo "$cC  cuttlebelle"
#   yarn global add cuttlebelle

  fnDivider

  fnInfo "Your public IP address in $cW $(curl --no-progress-meter https://api.ipify.org?format=text)$cX"

  fnDivider

fi

# ------------------------------------------------------------------------------
#   ___   _  _   ___   _____
#  |_ _| | \| | |_ _| |_   _|
#   | |  | .` |  | |    | |
#  |___| |_|\_| |___|   |_|
# ------------------------------------------------------------------------------

if [[ $ACTION == "init" ]] ; then

  echo -e "${Cgreen}gitpod$cC INIT$cT tasks"
  fnDivider

fi

# ------------------------------------------------------------------------------
#    ___    ___    __  __   __  __     _     _  _   ___
#   / __|  / _ \  |  \/  | |  \/  |   /_\   | \| | |   \
#  | (__  | (_) | | |\/| | | |\/| |  / _ \  | .` | | |) |
#   \___|  \___/  |_|  |_| |_|  |_| /_/ \_\ |_|\_| |___/
# ------------------------------------------------------------------------------

if [[ $ACTION == "command" ]] ; then

  echo -e "${Cgreen}gitpod$cC COMMAND$cT tasks"
  fnDivider

fi
