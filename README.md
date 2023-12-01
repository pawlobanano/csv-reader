# CSV file reader
CSV file reader with an e-mail's domain occurrences counter.

## Optimization research & ideas
1. Use buffered reading (*bufio package*),
2. Parallel processing (*processes email domains concurrently using worker goroutines*),
3. Optimize data structures (*structs for readability/maintainability improvement*),
4. Benchmark tests (*benchmarking different-sized input data files*),
5. Code profiling (*pprof tool to identify specific bottlenecks*).

## Environment variables
- To override config variables change the values in .env file. The default values:

    ```bash
    CONCURRENCY=4
    INPUT_CSV_FILE_PATH_DEFAULT=./data/test/customers_3k_lines.csv
    INPUT_CSV_FILE_PATH_0_LINES=../data/test/customers_0_lines.csv
    INPUT_CSV_FILE_PATH_10_LINES=../data/test/customers_10_lines.csv
    INPUT_CSV_FILE_PATH_3K_LINES=../data/test/customers_3k_lines.csv
    INPUT_CSV_FILE_PATH_10M_LINES=../data/test/customers_10m_lines.csv*
    READ_BUFFER_SIZE_IN_BYTES=4096
    ```
<sub>* _customers_10m_lines.csv_ file is stored locally due to the size (over 500 MB). It is used in benchmark tests.</sub>

## Screenshots from benchmark execution
- CONCURRENCY=1, READ_BUFFER_SIZE_IN_BYTES=4096
![CONCURRENCY=1, READ_BUFFER_SIZE_IN_BYTES=4096](/assets/benchmark-500ms-concurrency-1-read-buffer-size-4096-amd-ryzen-5-7600x.png)

- CONCURRENCY=6, READ_BUFFER_SIZE_IN_BYTES=4096
![CONCURRENCY=6, READ_BUFFER_SIZE_IN_BYTES=4096](/assets/benchmark-500ms-concurrency-6-read-buffer-size-4096-amd-ryzen-5-7600x.png)

- CONCURRENCY=12, READ_BUFFER_SIZE_IN_BYTES=4096
![CONCURRENCY=12, READ_BUFFER_SIZE_IN_BYTES=4096](/assets/benchmark-500ms-concurrency-12-read-buffer-size-4096-amd-ryzen-5-7600x.png)

- CONCURRENCY=1, READ_BUFFER_SIZE_IN_BYTES=8192
![CONCURRENCY=1, READ_BUFFER_SIZE_IN_BYTES=8192](/assets/benchmark-500ms-concurrency-1-read-buffer-size-8192-amd-ryzen-5-7600x.png)

- CONCURRENCY=6, READ_BUFFER_SIZE_IN_BYTES=8192
![CONCURRENCY=6, READ_BUFFER_SIZE_IN_BYTES=8192](/assets/benchmark-500ms-concurrency-6-read-buffer-size-8192-amd-ryzen-5-7600x.png)

- CONCURRENCY=12, READ_BUFFER_SIZE_IN_BYTES=8192
![CONCURRENCY=12, READ_BUFFER_SIZE_IN_BYTES=8192](/assets/benchmark-500ms-concurrency-12-read-buffer-size-8192-amd-ryzen-5-7600x.png)

- CONCURRENCY=1, READ_BUFFER_SIZE_IN_BYTES=16384
![CONCURRENCY=1, READ_BUFFER_SIZE_IN_BYTES=16384](/assets/benchmark-500ms-concurrency-1-read-buffer-size-16384-amd-ryzen-5-7600x.png)

- CONCURRENCY=6, READ_BUFFER_SIZE_IN_BYTES=16384
![CONCURRENCY=6, READ_BUFFER_SIZE_IN_BYTES=16384](/assets/benchmark-500ms-concurrency-6-read-buffer-size-16384-amd-ryzen-5-7600x.png)

- CONCURRENCY=12, READ_BUFFER_SIZE_IN_BYTES=16384
![CONCURRENCY=12, READ_BUFFER_SIZE_IN_BYTES=16384](/assets/benchmark-500ms-concurrency-12-read-buffer-size-16384-amd-ryzen-5-7600x.png)

## Makefile
- Run program
    ```
    make run
    ```

- Run tests
    ```
    make test
    ```

- Run benchmark
    ```
    make benchmark
    ```
