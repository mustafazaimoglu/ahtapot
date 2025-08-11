# 🐙 AHTAPOT - ScyllaDB/Cassandra Backup & Restore Tool

**AHTAPOT**, ScyllaDB ve Apache Cassandra veritabanları için **kolay, hızlı ve script tabanlı** bir **yedekleme (backup)** ve **geri yükleme (restore)** aracıdır.  
`dsbulk` ve `cqlsh` kullanarak **şema** (schema) ve **veri** (data) yedeklerini alır ve geri yükler.  
Küçük ve orta ölçekli ortamlar için idealdir.

---

## 🚀 Özellikler

- **Tam otomatik** backup & restore
- Tüm keyspace’ler, belirli bir keyspace veya tek bir tablo ile çalışma
- Şema ve/veya veri yedekleme seçeneği
- **json** veya **csv** formatında veri yedekleme
- **dsbulk** ve **cqlsh** tabanlı
- **Ahtapot imzalı yedek dosyası** ile geçerli backup kontrolü
- Yedekleme süresi ölçümü ve raporlama

---

## 📦 Gereksinimler

- **ScyllaDB** veya **Apache Cassandra**
- `cqlsh` erişimi
- [`dsbulk`](https://github.com/datastax/dsbulk) (DataStax Bulk Loader)
- Bash (Linux/Unix tabanlı sistemler)

---

## 🔧 Kurulum

1. Bu repo’yu klonlayın:
   ```bash
   git clone https://github.com/mustafazaimoglu/ahtapot.git
   cd ahtapot

   vi ahtapot.sh # dsbluk executable dosyasının yolu verin
   DSBULK=dsbulk-1.11.0/bin/dsbulk

   chmod +x ahtapot.sh
   ```
--- 


## 📂 Yedek Dosya Yapısı
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

## 📌 Kullanım
1. Tüm Keyspace’leri Yedekleme
    ```
    ./ahtapot.sh -h 127.0.0.1 -P 9042 -u cassandra -p cassandra \
    --all-keyspaces \
    --directory ./backup_2025_08_11 \
    --format json \
    --operation backup
    ```

2. Belirli Bir Keyspace’i Yedekleme
    ```
    ./ahtapot.sh -h 127.0.0.1 -P 9042 -u cassandra -p cassandra \
    --keyspace my_keyspace \
    --directory ./backup_myks \
    -f csv \
    -o backup
    ```

3. Belirli Bir Tabloyu Yedekleme
    ```
    ./ahtapot.sh -h 127.0.0.1 -P 9042 -u cassandra -p cassandra \
    --keyspace my_keyspace \
    --table my_table \
    --directory ./backup_table \
    -f json \
    -o backup
    ```

4. Yedekten Geri Yükleme (Restore)
    ```
    ./ahtapot.sh -h 127.0.0.1 -P 9042 -u cassandra -p cassandra \
    --all-keyspaces \
    --directory ./backup_2025_08_11 \
    -o restore
    ```

<br>
<p align="right">
<strong><i>Mustafa ZAİMOĞLU</i></strong>
</p>