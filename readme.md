# 🐙 AHTAPOT - ScyllaDB/Cassandra Backup & Restore Tool

**AHTAPOT** is an **easy, fast, and script-based** **backup** and **restore** tool for ScyllaDB and Apache Cassandra databases.  
It uses `dsbulk` and `cqlsh` to backup and restore **schema** and **data**.  
Ideal for small to medium-sized environments.

---

## 🚀 Features

- Fully automated backup & restore  
- Works with all keyspaces, a specific keyspace, or a single table  
- Option to backup schema and/or data  
- Data backup in **json** or **csv** format  
- Based on **dsbulk** and **cqlsh**  
- Valid backup verification with **Ahtapot signed backup file**  
- Backup duration measurement and reporting  

---

## 📦 Requirements

- **ScyllaDB** or **Apache Cassandra**  
- `cqlsh` 
- [`dsbulk`](https://github.com/datastax/dsbulk) (DataStax Bulk Loader) 
- Java (8 or later) 
- Python
- Bash (Linux/Unix based systems)  

---

## 🔧 Installation

1. Clone this repo
   ```bash
   git clone https://github.com/mustafazaimoglu/ahtapot.git
   cd ahtapot

   chmod +x ahtapot.sh
   ```

2. Install dsbulk (check the [official dsbulk releases](https://github.com/datastax/dsbulk/releases) for possible newer versions)
   ```bash
    wget https://github.com/datastax/dsbulk/releases/download/1.11.0/dsbulk-1.11.0.tar.gz
    tar -xzvf dsbulk-1.11.0.tar.gz
    mv dsbulk-1.11.0 /opt

    chmod +x /opt/dsbulk-1.11.0/bin/dsbulk

    ln -s /opt/dsbulk-1.11.0/bin/dsbulk /usr/bin/dsbulk
   ```

3. Install other requirements from above

--- 


## 📂 Backup File Structure
```bash
  [BACKUP_DIR]
  ├─ [KEYSPACE_NAME]
  │   ├─ keyspace.cql
  │   ├─ keyspace_full.cql
  │   └─ [TABLE_NAME]
  │       ├─ table.cql
  │       ├─ dump/
  │       └─ operation/
  └─ ahtapot
```
--- 

## 📌 Usage
1. Backup all keyspaces
    ```
    ./ahtapot.sh -h 127.0.0.1 -P 9042 -u cassandra -p cassandra \
    --all-keyspaces \
    --directory ./backup_2025_03_16 \
    --format json \
    --operation backup
    ```

2. Backup a specific keyspace
    ```
    ./ahtapot.sh -h 127.0.0.1 -P 9042 -u cassandra -p cassandra \
    -k my_keyspace \
    -d ./backup_myks \
    -f csv \
    -o backup
    ```

3. Backup a specific table
    ```
    ./ahtapot.sh -h 127.0.0.1 -P 9042 -u cassandra -p cassandra \
    -k my_keyspace \
    -t my_table \
    -d ./backup_table \
    -f json \
    -o backup
    ```

4. Restore all keyspaces from backup
    ```
    ./ahtapot.sh -h 127.0.0.1 -P 9042 -u cassandra -p cassandra \
    -A -d ./backup_2025_03_16 \
    -o restore
    ```

<br>
<p align="right">
<strong><i>Mustafa ZAİMOĞLU</i></strong>
</p>