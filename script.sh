slash () {
	TEMP_FILE="/tmp/slashnav-tomove.txt"
  ARGS=("$@")

  export SLASHNAV_WRAPPER=1
  slashnav "${ARGS[@]}"
  TOMOVE=$(cat "$TEMP_FILE")
  if [ -z "$TOMOVE" ]; then else
    cd "$TOMOVE"
    echo "" > "$TEMP_FILE"
  fi
  
}