# common file data for building opentsg
#                             _                       _
#   ___   _ __   ___   _ _   | |_   ___  __ _   ___  (_)  ___
#  / _ \ | '_ \ / -_) | ' \  |  _| (_-< / _` | |___| | | / _ \
#  \___/ | .__/ \___| |_||_|  \__| /__/ \__, |       |_| \___/
#        |_|                            |___/
#           common (+ link to docs)
#      ↗     ↗       ↑      ↖       ↖
#  check → build → test → deploy → postfix
#
#  opentsg-node component: github build / deploy / release
#
# --- clog functions ----------------------------------------------------------
[ -f clogrc/inc.sh ] && source clogrc/inc.sh
export PROJECT=$(basename $(pwd))        # generic project grabber
export isZSH="$ZSH_VERSION"              # [ -n "isZSH" ] && echo "shell=zsh"

# --- colors ------------------------------------------------------------------
c="\e[3"; R=1m; G=2m; B=4m; C=6m; M=5m; Y=3m; K=0m; W=7m; cX="\e[0m"
cC=$c$B; cU=$c$C; cT=$c$K; cI=$c$Y; cE=$c$R; cW=$c$M; cS=$c$G; cF=$c$W; cH=$c$C
# printf "$cC cmd$cU url$cT txt$cI inf$cE err$cW wrn$cS ok!$cF fle$cH hdr$cX\n"

# --- tagging ------------------------------------------------------------------

function fTagLocal  () {  git tag -a "$1" HEAD -m "$2" ; }
function fTagRemote () {  git push origin "$1" ; }

# --- user interaction --------------------------------------------------------

# --- comon data --------------------------------------------------------------

# fPrompt "Continue?" "yN" 7 ; index=$?       # 0 based index or len() after 7s
#
function fPrompt () {
  local SEC="$3";  [ -z "$SEC" ] && SEC=5        # timeout in seconds
  local OPT="$2";  [ -z "$OPT" ] && OPT="Yn"     # valid option strings
  local DFT=$(echo "$OPT" |egrep -o "[A-Z]")     # default response is in CAPS
  local opt=$( echo $OPT | tr '[:upper:]' '[:lower:]' )
  local charsBeforeDefaultChar=${OPT%%$DFT*}
  local defaultIndex=${#charsBeforeDefaultChar}

  printf -- "$1 ($OPT) "                         # print the prompt and options
  while true; do
    if [[ -n "$isZSH" ]] ; then
      read -t $SEC -k key                               # zsh read char in time
      err=$?
    else
      read -t $SEC -n 1 key                                    # bash read char
      err=$?
    fi
    [ $err -gt 0 ] &&  printf "\n" && return $defaultIndex     # handle timeout
    key=$( echo $key | tr '[:upper:]' '[:lower:]' )            # make lowercase
    if [[ $opt == *"$key"* ]] ; then                 # key is in the opt string
      prefix=${opt%%$key*}
      printf "\n"                                      # make output look right
      return ${#prefix}                     # return 0 based index or len($OPT)
    fi
    printf -- "\b \b"                      # erase char not in $opt & try again
  done
}

# --- end ---------------------------------------------------------------------
