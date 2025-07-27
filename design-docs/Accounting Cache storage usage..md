
##### Question
please calculate amount of data in redis with keys like this "balance:{user_id:int64}" and value of int64.
##### Answer
Let's calculate the **memory usage in Redis** for storing `100,000` keys of the form:

```
key:    balance:{user_id}
value: int64 (8 bytes)
```

---

### 🔢 1. Key Size Estimation

Each key looks like:  
`balance:{user_id}`  
Example: `balance:123456789`

#### Breakdown:

- `"balance:"` = 8 bytes
    
- `user_id` (int64, max 19 digits) ≈ up to 19 ASCII chars (usually fewer, but assume max)
    
- `:` = 1 byte
    

➡️ **Key length ≈ 8 + 1 + 19 = 28 bytes (worst case)**

> Redis keys are stored as **SDS (Simple Dynamic Strings)** which add some overhead (around 8–9 bytes for metadata).

So:  
➡️ **Key total size ≈ 28 bytes + 9 bytes overhead = ~37 bytes**

We’ll round up to **40 bytes** per key (safe estimate).

---

### 🔢 2. Value Size Estimation

The value is a **64-bit signed integer**, stored as a Redis **string**.

- Redis stores numbers as strings internally (e.g., `"1000"`), so:
    
    - `"9223372036854775807"` = max `int64` = 19 bytes
        
    - Average integer ≈ 8–12 ASCII bytes
        
- Add SDS overhead: 9 bytes
    

➡️ **Value total size ≈ 12 + 9 = ~21 bytes**  
We round up to **24 bytes** per value.

---

### 📦 3. Total Per-Key Size

|Component|Size (Bytes)|
|---|---|
|Key|~40|
|Value|~24|
|**Total**|**64 bytes**|

---

### 📊 4. Total Redis Memory for 100,000 Keys

```
64 bytes/key × 100,000 keys = 6,400,000 bytes ≈ 6.1 MB
```

✅ **Result: ~6.1 MB of RAM** to store 100,000 `int64` balances in Redis.

---

### 🔄 Optional Overheads (if you use):

|Feature|Overhead Estimate|
|---|---|
|Expiration TTL|+8–16 bytes/key|
|Hash instead of key-per-value|~half memory if keys are short|
|Replication|Memory x2|
|AOF persistence|Disk overhead only|

---

### 🧠 Optimization Tip

You can use a **Redis hash** like:

```text
HSET balances {user_id} {int64_value}
```

- One hash with 100,000 fields
    
- Much more space-efficient:
    
    - Keys are stored once (the hash key), not repeated
        
    - Uses **less metadata per field**
        

With a Redis hash, **total memory could drop to ~3–4 MB**.
