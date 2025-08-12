#!/bin/bash
echo ".....::::: AHTAPOT :::::....."

print_parameters() {
    echo "HOST: $HOST"
    echo "PORT: $PORT"
    echo "USERNAME: $USERNAME"
    echo "PASSWORD: $PASSWORD"
    echo "NO_SCHEMA: $NO_SCHEMA"
    echo "ALL_KEYSPACES: $ALL_KEYSPACES"
    echo "KEYSPACE: $KEYSPACE"
    echo "TABLE: $TABLE"
    echo "DIRECTORY: $DIRECTORY"
    echo "FORMAT: $FORMAT"
    echo "OPERATION: $OPERATION"
    echo "CONSISTENCY: $CONSISTENCY"
}

print_help() {
    echo "Usage:"
    echo "  $0 [parameters]"
    echo ""
    echo "Parameters:"
    echo "  --host, -h              Host address (default: 127.0.0.1)"
    echo "  --port, -P              Port number (default: 9042)"
    echo "  --username, -u          Username (default: cassandra)"
    echo "  --password, -p          Password (default: cassandra)"
    echo "  --no-schema, -S         Do not restore schema (Only restore data)"
    echo "  --no-data, -D           Do not backup/restore data (Only schema)"
    echo "  --all-keyspaces, -A     All keyspaces"
    echo "  --keyspace, -k          Keyspace name"
    echo "  --table, -t             Table name"
    echo "  --format, -f            Format (json | csv) (default: json)"
    echo "  --directory, -d         Directory path (target for backup, source for restore)"
    echo "  --operation, -o         Operation type (backup | restore) (required)"
    echo "  --consistency, -c       Consistency level (ANY | LOCAL_ONE | ONE | TWO | THREE | LOCAL_QUORUM | QUORUM | EACH_QUORUM | ALL) (default: LOCAL_ONE)"
    echo "  --help, -H              Help"
    echo ""
}

check_connection(){
    local output
    if ! output=$(cqlsh "$HOST" "$PORT" -u "$USERNAME" -p "$PASSWORD" -e "SELECT version FROM system.versions LIMIT 1"); then
        echo "ERROR > Connection failed. Please check the connection details." >&2
        exit 1
    fi
}

get_existing_keyspaces_from_db() {
  local -n result_array="$1"  # array referansı (nameref)

  readarray -t result_array < <(
    cqlsh $HOST $PORT -u $USERNAME -p $PASSWORD -e "SELECT keyspace_name FROM system_schema.keyspaces" 2>/dev/null \
    | tail -n +4 | head -n -2 \
    | awk '{$1=$1; print}' \
    | grep -v -E '^(system|system_schema|system_auth|system_traces|system_distributed_everywhere|system_distributed)$' \
    | sort
  )
}

get_tables_in_keyspace_from_db() {
  local keyspace="$1"
  local -n result_array="$2"  # 'nameref' ile çağırılan array'e doğrudan yaz

  readarray -t result_array < <(
    cqlsh $HOST $PORT -u $USERNAME -p $PASSWORD -e "SELECT table_name FROM system_schema.tables WHERE keyspace_name = '${keyspace}'" 2>/dev/null \
    | tail -n +4 | head -n -2 \
    | awk '{$1=$1; print}' \
    | sort
  )
}

get_directories_from_dir() {
    local main_dir="$1"
    local -n result_array="$2"  # nameref: array referansı

    mapfile -t result_array < <(find "$main_dir" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
}

create_ahtapot_file() {
    local main_dir="$1"
    local file="$main_dir/ahtapot"

    mkdir -p $main_dir
    touch $file
    echo $FORMAT > $file
    echo "Coded By MZ" >> $file
}

create_directories_for_keyspace() {
    local main_dir="$1"
    local keyspace_name="$2"

    mkdir -p $main_dir/$keyspace_name
}

create_directories_for_table() {
    local main_dir="$1"
    local keyspace_name="$2"
    local table_name="$3"
    
    mkdir -p $main_dir/$keyspace_name/$table_name/dump
    mkdir -p $main_dir/$keyspace_name/$table_name/operation
}

backup_schema_keyspace() {
    local main_dir="$1"
    local keyspace_name="$2"
    local full_path="$main_dir/$keyspace_name"

    echo ">>> BACKUP KEYSPACE SCHEMA : $keyspace_name"
    cqlsh $HOST $PORT -u $USERNAME -p $PASSWORD -e "DESC KEYSPACE \"$keyspace_name\"" 2>/dev/null | head -n 3 > $full_path/keyspace.cql 
    cqlsh $HOST $PORT -u $USERNAME -p $PASSWORD -e "DESC KEYSPACE \"$keyspace_name\"" 2>/dev/null > $full_path/keyspace_full.cql 
}

backup_schema_table() {
    local main_dir="$1"
    local keyspace_name="$2"
    local table_name="$3"

    local full_path="$main_dir/$keyspace_name/$table_name"

    echo ">>> BACKUP TABLE SCHEMA: $keyspace_name.$table_name"
    cqlsh $HOST $PORT -u $USERNAME -p $PASSWORD -e "DESC TABLE \"$keyspace_name\".\"$table_name\"" > $full_path/table.cql 2>/dev/null
}

backup_table() {
    local main_dir="$1"
    local keyspace_name="$2"
    local table_name="$3"

    local full_path="$main_dir/$keyspace_name/$table_name"

    local backup_path="$full_path/dump"
    local log_path="$full_path/operation"

    echo ">>> BACKUP TABLE DATA: $keyspace_name.$table_name"
    dsbulk unload -u $USERNAME -p $PASSWORD -h $HOST -port $PORT -cl $CONSISTENCY -k $keyspace_name -t $table_name -url $backup_path -c $FORMAT -logDir $log_path
    echo ""
}

restore_schema_keyspace() {
    local main_dir="$1"
    local keyspace_name="$2"
    
    local ddl_file="$main_dir/$keyspace_name/keyspace.cql"

    echo ">>> RESTORE KEYSPACE SCHEMA : $keyspace_name"
    cqlsh $HOST $PORT -u $USERNAME -p $PASSWORD -f $ddl_file
}

restore_schema_table() {
    local main_dir="$1"
    local keyspace_name="$2"
    local table_name="$3"
    
    local ddl_file="$main_dir/$keyspace_name/$table_name/table.cql"

    echo ">>> RESTORE TABLE SCHEMA: $keyspace_name.$table_name"
    cqlsh $HOST $PORT -u $USERNAME -p $PASSWORD -f $ddl_file
}

restore_table() {
    local main_dir="$1"
    local keyspace_name="$2"
    local table_name="$3"

    local backup_path="$main_dir/$keyspace_name/$table_name/dump"
    local log_path=""$main_dir"_restore_logs/$keyspace_name/$table_name"
    mkdir -p $log_path

    echo ">>> RESTORE TABLE DATA: $keyspace_name.$table_name"
    dsbulk load -u $USERNAME -p $PASSWORD -h $HOST -port $PORT -cl $CONSISTENCY -k $keyspace_name -t $table_name -url $backup_path -c $FORMAT -logDir $log_path
    echo ""
}

# FILE STRUCTURE:
# directory_name
#   keyspace_name
#     keyspace.cql
#     keyspace_full.cql
#     table_name 
#       - table.cql
#       - dump
#       - operation

# Varsayılanlar
HOST="127.0.0.1"
PORT="9042"
USERNAME="cassandra"
PASSWORD="cassandra"
NO_SCHEMA="false"
NO_DATA="false"
FORMAT="json"
CONSISTENCY="LOCAL_ONE"
START_TIME=$(date +%s)

# Parametreleri oku
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --host|-h) HOST="$2"; shift ;;
        --port|-P) PORT="$2"; shift ;;
        --username|-u) USERNAME="$2"; shift ;;
        --password|-p) PASSWORD="$2"; shift ;;
        --no-schema|-S) NO_SCHEMA=true; ;;
        --no-data|-D) NO_DATA=true; ;;
        --all-keyspaces|-A) ALL_KEYSPACES=true; ;;
        --keyspace|-k) KEYSPACE="$2"; shift ;;
        --table|-t) TABLE="$2"; shift ;;
        --directory|-d) DIRECTORY="$2"; shift ;;
        --format|-f) FORMAT="$2"; shift ;;
        --operation|-o) OPERATION="$2"; shift ;;
        --consistency|-c) CONSISTENCY="$2"; shift ;;
        --help|-H) print_help; exit 0 ;;
        *) echo "❌ Unknown parameter: $1"; print_help; exit 1 ;;
    esac
    shift
done

# print parameters
# print_parameters

# cqlsh check
if ! command -v cqlsh > /dev/null 2>&1; then
    echo "ERROR > 'cqlsh' not found in PATH. Please install it or update your PATH."
    exit 1
fi

# dsbulk check
if ! command -v dsbulk > /dev/null 2>&1; then
    echo "ERROR > 'dsbulk' not found in PATH. Please install it or update your PATH."
    exit 1
fi

# --operation zorunlu
if [[ -z "$OPERATION" ]]; then
    echo "ERROR > Operation type is mandatory!"
    exit 1
fi

# --directory zorunlu
if [[ -z "$DIRECTORY" ]]; then
    echo "ERROR > Backup/Restore directory is mandatory!"
    exit 1
fi

# mode control
OPERATION_MODE=0;
# 1 = ALL KEYSPACES, 2 = SPECIFIC KEYSPACE, 3 = SPECIFIC TABLE 
if [[ "$ALL_KEYSPACES" == true ]]; then
    if [[ -n "$KEYSPACE" || -n "$TABLE" ]]; then
        echo "ERROR > --keyspace or --table cannot be used together with --all-keyspaces."
        exit 1
    fi
    OPERATION_MODE=1;
elif [[ -n "$KEYSPACE" && -z "$TABLE" ]]; then
    OPERATION_MODE=2;
elif [[ -n "$KEYSPACE" && -n "$TABLE" ]]; then
    OPERATION_MODE=3;
else
    echo "ERROR > Invalid parameter combination."
    echo "Only one of the following is allowed:"
    echo "  1) --all-keyspaces                    # For all keyspaces"
    echo "  2) --keyspace KEYSPACE                # All tables in the specified keyspace"
    echo "  3) --keyspace KEYSPACE --table TABLE  # A specific table"
    exit 1
fi

# --consistency check
CONSISTENCY=$(echo "$CONSISTENCY" | tr '[:lower:]' '[:upper:]')
case "$CONSISTENCY" in
  ANY|LOCAL_ONE|ONE|TWO|THREE|LOCAL_QUORUM|QUORUM|EACH_QUORUM|ALL)
    :
    ;;
  *)
    echo "ERROR > Invalid consistency level!"
    echo "Valid consistency levels are: (ANY | LOCAL_ONE | ONE | TWO | THREE | LOCAL_QUORUM | QUORUM | EACH_QUORUM | ALL)"
    exit 1
    ;;
esac

# İşleme göre kontrol ve komut üretimi
if [[ "$OPERATION" == "backup" ]]; then
    # --directory existence check
    if [[ -d $DIRECTORY ]]; then
        echo "ERROR > Directory exists!"
        exit 1
    fi

    # --format check
    case "$FORMAT" in
    json|csv)
        :
        ;;
    *)
        echo "ERROR > Invalid format type!"
        echo "Valid format types are: (json | csv)"
        exit 1
        ;;
    esac

    if [[ "$NO_SCHEMA" == "true" ]]; then
        echo "WARNING > -S,--no-schema parameter ignored for operation backup!"
    fi

    # CHECK CONNECTION
    check_connection

    # BACKUP LOGIC
    create_ahtapot_file $DIRECTORY 

    declare -a EXISTING_KEYSPACES
    get_existing_keyspaces_from_db EXISTING_KEYSPACES

    if [[ "$OPERATION_MODE" -eq 1 ]]; then
        for ks in "${EXISTING_KEYSPACES[@]}"; do
            declare -a EXISTING_TABLES_IN_KEYSPACE
            get_tables_in_keyspace_from_db "$ks" EXISTING_TABLES_IN_KEYSPACE

            create_directories_for_keyspace $DIRECTORY $ks
            backup_schema_keyspace $DIRECTORY $ks

            for tb in "${EXISTING_TABLES_IN_KEYSPACE[@]}"; do
                create_directories_for_table $DIRECTORY $ks $tb
                backup_schema_table $DIRECTORY $ks $tb

                if [[ "$NO_DATA" == "false" ]]; then
                    backup_table $DIRECTORY $ks $tb
                fi
            done
        done
    elif [[ "$OPERATION_MODE" -eq 2 ]]; then
        if [[ ! " ${EXISTING_KEYSPACES[@]} " =~ " ${KEYSPACE} " ]]; then
            echo "ERROR > Keyspace '$KEYSPACE' doesn't exists!"
            exit 1
        fi

        declare -a EXISTING_TABLES_IN_KEYSPACE
        get_tables_in_keyspace_from_db "$KEYSPACE" EXISTING_TABLES_IN_KEYSPACE

        create_directories_for_keyspace $DIRECTORY $KEYSPACE
        backup_schema_keyspace $DIRECTORY $KEYSPACE

        for tb in "${EXISTING_TABLES_IN_KEYSPACE[@]}"; do
            create_directories_for_table $DIRECTORY $KEYSPACE $tb
            backup_schema_table $DIRECTORY $KEYSPACE $tb

            if [[ "$NO_DATA" == "false" ]]; then
                backup_table $DIRECTORY $KEYSPACE $tb
            fi
        done
    elif [[ "$OPERATION_MODE" -eq 3 ]]; then
        if [[ ! " ${EXISTING_KEYSPACES[@]} " =~ " ${KEYSPACE} " ]]; then
            echo "ERROR > Keyspace '$KEYSPACE' doesn't exists!"
            exit 1
        fi

        declare -a EXISTING_TABLES_IN_KEYSPACE
        get_tables_in_keyspace_from_db "$KEYSPACE" EXISTING_TABLES_IN_KEYSPACE

        if [[ ! " ${EXISTING_TABLES_IN_KEYSPACE[@]} " =~ " ${TABLE} " ]]; then
            echo "ERROR > Table '$KEYSPACE'.'$TABLE' doesn't exists!"
            exit 1
        fi
        
        create_directories_for_keyspace $DIRECTORY $KEYSPACE
        backup_schema_keyspace $DIRECTORY $KEYSPACE

        create_directories_for_table $DIRECTORY $KEYSPACE $TABLE
        backup_schema_table $DIRECTORY $KEYSPACE $TABLE

        if [[ "$NO_DATA" == "false" ]]; then
            backup_table $DIRECTORY $KEYSPACE $TABLE
        fi
    fi
elif [[ "$OPERATION" == "restore" ]]; then
    if [[ ! -d "$DIRECTORY" ]]; then
        echo "ERROR > Directory does not exist!"
        exit 1
    elif [[ -z "$(ls -A "$DIRECTORY" 2>/dev/null)" ]]; then
        echo "ERROR > Directory is empty!"
        exit 1
    elif [[ ! -f "$DIRECTORY/ahtapot" ]]; then
        echo "ERROR > This is not a valid backup!"
        exit 1
    fi

    if [[ "$NO_DATA" == true && "$NO_SCHEMA" == true ]]; then
        echo "ERROR > -D,--no-data and -S,--no-schema cannot be used together!"
        exit 1
    fi

    # CHECK CONNECTION
    check_connection

    # RESTORE LOGIC
    FORMAT=$(head -n 1 $DIRECTORY/ahtapot)
    echo "Restore format set to '$FORMAT' because backup was taken in '$FORMAT' format!"

    declare -a EXISTING_KEYSPACES_RESTORE
    get_directories_from_dir $DIRECTORY EXISTING_KEYSPACES_RESTORE

    if [[ "$OPERATION_MODE" -eq 1 ]]; then
        for ksr in "${EXISTING_KEYSPACES_RESTORE[@]}"; do
            declare -a EXISTING_TABLES_IN_KEYSPACE_RESTORE
            get_directories_from_dir "$DIRECTORY/$ksr" EXISTING_TABLES_IN_KEYSPACE_RESTORE

            if [[ "$NO_SCHEMA" == "false" ]]; then
                restore_schema_keyspace $DIRECTORY $ksr
            fi

            for tbr in "${EXISTING_TABLES_IN_KEYSPACE_RESTORE[@]}"; do
                if [[ "$NO_SCHEMA" == "false" ]]; then
                    restore_schema_table $DIRECTORY $ksr $tbr
                fi

                if [[ "$NO_DATA" == "false" ]]; then
                    restore_table $DIRECTORY $ksr $tbr
                fi
            done
        done
    elif [[ "$OPERATION_MODE" -eq 2 ]]; then
        if [[ ! " ${EXISTING_KEYSPACES_RESTORE[@]} " =~ " ${KEYSPACE} " ]]; then
            echo "ERROR > Keyspace '$KEYSPACE' doesn't exists in the backup!"
            exit 1
        fi

        declare -a EXISTING_TABLES_IN_KEYSPACE_RESTORE
        get_directories_from_dir "$DIRECTORY/$KEYSPACE" EXISTING_TABLES_IN_KEYSPACE_RESTORE

        if [[ "$NO_SCHEMA" == "false" ]]; then
            restore_schema_keyspace $DIRECTORY $KEYSPACE
        fi

        for tbr in "${EXISTING_TABLES_IN_KEYSPACE_RESTORE[@]}"; do
            if [[ "$NO_SCHEMA" == "false" ]]; then
                restore_schema_table $DIRECTORY $KEYSPACE $tbr
            fi
            
            if [[ "$NO_DATA" == "false" ]]; then
                restore_table $DIRECTORY $KEYSPACE $tbr
            fi
        done
    elif [[ "$OPERATION_MODE" -eq 3 ]]; then
        if [[ ! " ${EXISTING_KEYSPACES_RESTORE[@]} " =~ " ${KEYSPACE} " ]]; then
            echo "ERROR > Keyspace '$KEYSPACE' doesn't exists in the backup!"
            exit 1
        fi

        declare -a EXISTING_TABLES_IN_KEYSPACE_RESTORE
        get_directories_from_dir "$DIRECTORY/$KEYSPACE" EXISTING_TABLES_IN_KEYSPACE_RESTORE

        if [[ ! " ${EXISTING_TABLES_IN_KEYSPACE_RESTORE[@]} " =~ " ${TABLE} " ]]; then
            echo "ERROR > Table '$KEYSPACE'.'$TABLE' doesn't exists in the backup!"
            exit 1
        fi

        if [[ "$NO_SCHEMA" == "false" ]]; then
            restore_schema_keyspace $DIRECTORY $KEYSPACE
            restore_schema_table $DIRECTORY $KEYSPACE $TABLE
        fi
        
        if [[ "$NO_DATA" == "false" ]]; then
            restore_table $DIRECTORY $KEYSPACE $TABLE
        fi
    fi
else
    echo "ERROR > Operation type must be either 'backup' or 'restore'."
    print_help
    exit 1
fi

END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

# Süreyi saat:dakika:saniye formatına çevir
HOURS=$((DURATION / 3600))
MINUTES=$(( (DURATION % 3600) / 60 ))
SECONDS=$((DURATION % 60))

# Formatlı çıktı
printf -v FORMATTED "%02d:%02d:%02d" "$HOURS" "$MINUTES" "$SECONDS" 
echo -e "\e[32m${OPERATION^^} Completed!\e[0m -> \e[36mDuration: $FORMATTED\e[0m"

if [[ "$OPERATION" == "backup" ]]; then
    echo "Duration: $FORMATTED" >> $DIRECTORY/ahtapot
fi
#Coded By MZ