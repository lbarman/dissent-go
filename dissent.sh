#!/usr/bin/env bash


# ************************************
# Dissent all-in-one startup script
# ************************************
# author : Ludovic Barman
# email : ludovic.barman@gmail.com
# belongs to : the PriFi project
# 			<github.com/lbarman/dissent-go>
# ************************************

# variables that you might change often

dbg_lvl=3                       # 1=less verbose, 3=more verbose. goes up to 5, but then prints the SDA's message (network framework)
try_use_real_identities="false"  # if "true", will try to use "self-generated" public/private key as a replacement for the dummy keys
                                # we generated for you. It asks you if it does not find real keys. If false, will always use the dummy keys.
colors="true"                   # if  "false", the output of PriFi (and this script) will be in black-n-white

# default file names :

prifi_file="dissent.toml"                     # default name for the prifi config file (contains prifi-specific settings)
identity_file="identity.toml"               # default name for the identity file (contains public + private key)
group_file="group.toml"                     # default name for the group file (contains public keys + address of other nodes)

# location of the buildable (go build) prifi file :

bin_file="$GOPATH/src/github.com/lbarman/dissent-go/app/dissent.go"

# we have two "identities" directory. The second one is empty unless you generate your own keys with "gen-id"

configdir="config"
defaultIdentitiesDir="identities_default"   # in $configdir
realIdentitiesDir="identities_real"         # in $configdir

sleeptime_between_spawns=1                  # time in second between entities launch in all-localhost part

source "helpers.lib.sh"

# ------------------------
#     HELPER FUNCTIONS
# ------------------------

print_usage() {
	echo
	echo -e "Dissent Simulation in go"
	echo
}

# ------------------------
#     MAIN SWITCH
# ------------------------

# $1 is operation : "install", "client", "trustee", "clean", "gen-id"
case $1 in

	install|Install|INSTALL)

		echo -n "Testing for GO... "
		test_go
		echo -e "$okMsg"

		echo -n "Getting all go packages... "
		cd app; go get ./... 1>/dev/null 2>&1
		cd ..
		echo -e "$okMsg"

		echo -en "Switching ONet branch to ${highlightOn}$cothorityBranchRequired${highlightOff}... "
		cd "$GOPATH/src/gopkg.in/dedis/onet.v2"; git checkout "$cothorityBranchRequired" 1>/dev/null 2>&1
		echo -e "$okMsg"

		echo -n "Re-getting all go packages (since we switched branch)... "
		cd "$GOPATH/src/github.com/lbarman/dissent-go/app"; go get ./... 1>/dev/null 2>&1
		cd ../..
		cd "$GOPATH/src/gopkg.in/dedis/onet.v2"; go get -u ./... 1>/dev/null 2>&1
		echo -e "$okMsg"

		echo -n "Testing ONet branch... "
		test_cothority
		echo -e "$okMsg"

		;;

	trustee|Trustee|TRUSTEE)

		trusteeId="$2"

		#test for proper setup
		test_go
		test_cothority

		if [ "$#" -lt 2 ]; then
			echo -e "$errorMsg parameter 2 need to be the trustee id."
			exit 1
		fi
		test_digit "$trusteeId" 2

		#specialize the config file (we use the dummy folder, and maybe we replace with the real folder after)
		prifi_file2="$configdir/$prifi_file"
		identity_file2="$configdir/$defaultIdentitiesDir/trustee$trusteeId/$identity_file"
		group_file2="$configdir/$defaultIdentitiesDir/trustee$trusteeId/$group_file"

		#we we want to, try to replace with the real folder
		if [ "$try_use_real_identities" = "true" ]; then
			if [ -f "$configdir/$realIdentitiesDir/trustee$trusteeId/$identity_file" ] && [ -f "$configdir/$defaultIdentitiesDir/trustee$trusteeId/$group_file" ]; then
				echo -e "$okMsg Found real identities (in $configdir/$realIdentitiesDir/trustee$trusteeId/), using those."
				identity_file2="$configdir/$realIdentitiesDir/trustee$trusteeId/$identity_file"
				group_file2="$configdir/$realIdentitiesDir/trustee$trusteeId/$group_file"
			else
				echo -e "$warningMsg Trying to use real identities, but does not exists for trustee $trusteeId (in $configdir/$realIdentitiesDir/trustee$trusteeId/). Falling back to pre-generated ones."
			fi
		else
			echo -e "$warningMsg using pre-created identities. Set \"try_use_real_identities\" to True in real deployements."
		fi

		# test that all files exists
		test_files

		#run PriFi in relay mode
		DEBUG_COLOR="$colors" go run "$bin_file" --cothority_config "$identity_file2" --group "$group_file2" -d "$dbg_lvl" --prifi_config "$prifi_file2" trustee
		;;

	client0|Client0|CLIENT0)

		#test for proper setup
		test_go
		test_cothority
		clientId=0

		#specialize the config file (we use the dummy folder, and maybe we replace with the real folder after)
		prifi_file2="$configdir/$prifi_file"
		identity_file2="$configdir/$defaultIdentitiesDir/client$clientId/$identity_file"
		group_file2="$configdir/$defaultIdentitiesDir/client$clientId/$group_file"

		#we we want to, try to replace with the real folder
		if [ "$try_use_real_identities" = "true" ]; then
			if [ -f "$configdir/$realIdentitiesDir/client$clientId/$identity_file" ] && [ -f "$configdir/$realIdentitiesDir/client$clientId/$group_file" ]; then
				echo -e "$okMsg Found real identities (in $configdir/$realIdentitiesDir/client$clientId/), using those."
				identity_file2="$configdir/$realIdentitiesDir/client$clientId/$identity_file"
				group_file2="$configdir/$realIdentitiesDir/client$clientId/$group_file"
			else
				echo -e "$warningMsg Trying to use real identities, but does not exists for client $clientId (in $configdir/$realIdentitiesDir/client$clientId/). Falling back to pre-generated ones."
			fi
		else
			echo -e "$warningMsg using pre-created identities. Set \"try_use_real_identities\" to True in real deployements."
		fi

		# test that all files exists
		test_files

		#run PriFi in relay mode
		DEBUG_COLOR="$colors" go run "$bin_file" --cothority_config "$identity_file2" --group "$group_file2" -d "$dbg_lvl" --prifi_config "$prifi_file2" client0
		;;

	client|Client|CLIENT)

		clientId="$2"

		#test for proper setup
		test_go
		test_cothority

		if [ "$#" -lt 2 ]; then
			echo -e "$errorMsg parameter 2 need to be the client id."
			exit 1
		fi
		test_digit "$clientId" 2

		# the 3rd argument can replace the port number
		if [ "$#" -eq 3 ]; then
			test_digit "$3" 3
			socksServer1Port="$3"
		fi

		#specialize the config file (we use the dummy folder, and maybe we replace with the real folder after)
		prifi_file2="$configdir/$prifi_file"
		identity_file2="$configdir/$defaultIdentitiesDir/client$clientId/$identity_file"
		group_file2="$configdir/$defaultIdentitiesDir/client$clientId/$group_file"

		#we we want to, try to replace with the real folder
		if [ "$try_use_real_identities" = "true" ]; then
			if [ -f "$configdir/$realIdentitiesDir/client$clientId/$identity_file" ] && [ -f "$configdir/$realIdentitiesDir/client$clientId/$group_file" ]; then
				echo -e "$okMsg Found real identities (in $configdir/$realIdentitiesDir/client$clientId/), using those."
				identity_file2="$configdir/$realIdentitiesDir/client$clientId/$identity_file"
				group_file2="$configdir/$realIdentitiesDir/client$clientId/$group_file"
			else
				echo -e "$warningMsg Trying to use real identities, but does not exists for client $clientId (in $configdir/$realIdentitiesDir/client$clientId/). Falling back to pre-generated ones."
			fi
		else
			echo -e "$warningMsg using pre-created identities. Set \"try_use_real_identities\" to True in real deployements."
		fi

		# test that all files exists
		test_files

		#run PriFi in relay mode
		DEBUG_COLOR="$colors" go run "$bin_file" --cothority_config "$identity_file2" --group "$group_file2" -d "$dbg_lvl" --prifi_config "$prifi_file2" client
		;;

	gen-id|Gen-Id|GEN-ID)
		echo -e "Going to generate private/public keys (named ${highlightOn}identity.toml${highlightOff})..."

		read -p "Do you want to generate it for [r]elay, [c]lient, or [t]trustee ? " key

		path=""
		case "$key" in
			r|R)
				path="relay"
			;;
			t|T)

				read -p "Do you want to generate it for trustee [0] or [1] (or more - enter digit) ? " key2

				test_digit "$key2" 1
				pathSource="trustee0"
				path="trustee$key2"
				;;

			c|C)
				read -p "Do you want to generate it for client [0],[1] or [2] (or more - enter digit) ? " key2

				test_digit "$key2" 1
				pathSource="client0"
				path="client$key2"
				;;


			*)
				echo -e "$errorMsg did not understand."
				exit 1
				;;
		esac

		pathReal="$configdir/$realIdentitiesDir/$path/"
		pathDefault="$configdir/$defaultIdentitiesDir/$pathSource/"
		echo -e "Gonna generate ${highlightOn}identity.toml${highlightOff} in ${highlightOn}$pathReal${highlightOff}"

		#generate identity.toml
		DEBUG_COLOR="$colors" go run "$bin_file" --default_path "$pathReal" gen-id

		if [ ! -f "${pathReal}group.toml" ]; then
			#now group.toml
			echo -n "Done ! now copying group.toml from identities_default/ to identity_real/..."
			cp "${pathDefault}/group.toml" "${pathReal}group.toml"
			echo -e "$okMsg"

			echo -e "Please edit ${highlightOn}$pathReal/group.toml${highlightOff} to the correct values."
		else
			echo -e "Group file ${highlightOn}$pathReal/group.toml${highlightOff} already exists, not overwriting! you might want to check that the contents are correct."
		fi
		;;

	clean|Clean|CLEAN)
		echo -n "Cleaning local log files... 			"
		rm *.log 1>/dev/null 2>&1
		echo -e "$okMsg"
		;;

	*)
		print_usage
		;;
esac
