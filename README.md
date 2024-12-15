# Smoke Test

Deep inside Initrode Global's enterprise management framework lies a component that writes data to a server and expects to read the same data back. (Think of it as a kind of distributed system [delay-line memory](https://en.wikipedia.org/wiki/Delay-line_memory)). We need you to write the server to echo the data back.

Accept TCP connections.

Whenever you receive data from a client, send it back unmodified.

Make sure you don't mangle binary data, and that you can handle at least 5 simultaneous clients.

Once the client has finished sending data to you it shuts down its sending side. Once you've reached end-of-file on your receiving side, and sent back all the data you've received, close the socket so that the client knows you've finished. (This point trips up a lot of proxy software, such as ngrok; if you're using a proxy and you can't work out why you're failing the check, try hosting your server in the cloud instead).

Your program will implement the TCP Echo Service from [RFC 862](https://www.rfc-editor.org/rfc/rfc862.html).

# Prime Time

To keep costs down, a hot new government department is contracting out its mission-critical primality testing to the lowest bidder. (That's you).

Officials have devised a [JSON](https://en.wikipedia.org/wiki/JSON)\-based request-response protocol. Each request is a single line containing a JSON object, terminated by a newline character (`'\n'`, or ASCII 10). Each request begets a response, which is also a single line containing a JSON object, terminated by a newline character.

After connecting, a client may send multiple requests in a single session. Each request should be handled in order.

A conforming request object has the required field `method`, which must always contain the string `"isPrime"`, and the required field `number`, which must contain a number. Any JSON number is a valid number, including floating-point values.

Example request:

    {"method":"isPrime","number":123}
    

A request is _malformed_ if it is not a well-formed JSON object, if any required field is missing, if the method name is not `"isPrime"`, or if the `number` value is not a number.

Extraneous fields are to be ignored.

A conforming response object has the required field `method`, which must always contain the string `"isPrime"`, and the required field `prime`, which must contain a boolean value: `true` if the number in the request was prime, `false` if it was not.

Example response:

    {"method":"isPrime","prime":false}
    

A response is _malformed_ if it is not a well-formed JSON object, if any required field is missing, if the method name is not `"isPrime"`, or if the `prime` value is not a boolean.

A response object is considered _incorrect_ if it is well-formed but has an incorrect `prime` value. Note that non-integers can not be prime.

Accept TCP connections.

Whenever you receive a conforming request, send back a correct response, and wait for another request.

Whenever you receive a _malformed_ request, send back a single _malformed_ response, and disconnect the client.

Make sure you can handle at least 5 simultaneous clients.

# Means to an end

Your friendly neighbourhood investment bank is having trouble analysing historical price data. They need you to build a TCP server that will let clients **insert and query timestamped prices**.

Overview
--------

Clients will connect to your server using TCP. Each client tracks the price of a **different asset**. Clients send messages to the server that either **insert** or **query** the prices.

Each connection from a client is a separate session. Each session's data represents a different asset, so each session can **only query the data supplied by itself**.

Message format
--------------

To keep bandwidth usage down, a simple **binary format** has been specified.

Each message from a client is **9 bytes** long. Clients can send multiple messages per connection. Messages are _not_ delimited by newlines or any other character: you'll know where one message ends and the next starts because they are always 9 bytes.

    Byte:  |  0  |  1     2     3     4  |  5     6     7     8  |
    Type:  |char |         int32         |         int32         |
    

The first byte of a message is a character indicating its _type_. This will be an [ASCII](https://en.wikipedia.org/wiki/ASCII) uppercase **`'I'`** or **`'Q'`** character, indicating whether the message _inserts_ or _queries_ prices, respectively.

The next 8 bytes are **two signed [two's complement](https://en.wikipedia.org/wiki/Two%27s_complement) 32-bit integers in network byte order** ([big endian](https://en.wikipedia.org/wiki/Endianness)), whose meaning depends on the message type. We'll refer to these numbers as **`int32`**, but note this may differ from your system's native `int32` type (if any), particularly with regard to byte order.

Behaviour is undefined if the type specifier is not either `'I'` or `'Q'`.

### Insert

An _insert_ message lets the client **insert a timestamped price**.

The message format is:

    Byte:  |  0  |  1     2     3     4  |  5     6     7     8  |
    Type:  |char |         int32         |         int32         |
    Value: | 'I' |       timestamp       |         price         |
    

The first `int32` is the _timestamp_, in seconds since [00:00, 1st Jan 1970](https://en.wikipedia.org/wiki/Unix_time).

The second `int32` is the _price_, in pennies, of this client's asset, at the given timestamp.

Note that:

*   Insertions _may_ occur out-of-order.
*   While rare, prices [can go negative](https://www.bbc.co.uk/news/business-52350082).
*   Behaviour is undefined if there are multiple prices with the same timestamp from the same client.

For example, to insert a price of _101_ pence at timestamp _12345_, a client would send:

    Hexadecimal: 49    00 00 30 39    00 00 00 65
    Decoded:      I          12345            101
    

(Remember that you'll receive 9 raw bytes, rather than ASCII text representing hex-encoded data).

### Query

A _query_ message lets the client **query the average price over a given time period**.

The message format is:

    Byte:  |  0  |  1     2     3     4  |  5     6     7     8  |
    Type:  |char |         int32         |         int32         |
    Value: | 'Q' |        mintime        |        maxtime        |
    

The first `int32` is _mintime_, the earliest timestamp of the period.

The second `int32` is _maxtime_, the latest timestamp of the period.

The server must compute the **mean** of the inserted prices with timestamps _T_, _**mintime <= T <= maxtime**_ (i.e. timestamps in the closed interval _\[mintime, maxtime\]_). If the mean is not an integer, it is acceptable to round **either up or down**, at the server's discretion.

The server must then **send the mean to the client** as a single `int32`.

If there are no samples within the requested period, or if _mintime_ comes after _maxtime_, the value returned must be 0.

For example, to query the mean price between _T=1000_ and _T=100000_, a client would send:

    Hexadecimal: 51    00 00 03 e8    00 01 86 a0
    Decoded:      Q           1000         100000
    

And if the mean price during this time period were _5107_ pence, the server would respond:

    Hexadecimal: 00 00 13 f3
    Decoded:            5107
    

(Remember that you'll receive 9 raw bytes, and send 4 raw bytes, rather than ASCII text representing hex-encoded data).

Example session
---------------

In this example, "`-->`" denotes messages from the server to the client, and "`<--`" denotes messages from the client to the server.

        Hexadecimal:                 Decoded:
    <-- 49 00 00 30 39 00 00 00 65   I 12345 101
    <-- 49 00 00 30 3a 00 00 00 66   I 12346 102
    <-- 49 00 00 30 3b 00 00 00 64   I 12347 100
    <-- 49 00 00 a0 00 00 00 00 05   I 40960 5
    <-- 51 00 00 30 00 00 00 40 00   Q 12288 16384
    --> 00 00 00 65                  101
    

The client inserts _(timestamp,price)_ values: _(12345,101)_, _(12346,102)_, _(12347,100)_, and _(40960,5)_. The client then queries between _T=12288_ and _T=16384_. The server computes the mean price during this period, which is _101_, and sends back _101_.

Other requirements
------------------

Make sure you can handle at least 5 simultaneous clients.

Where a client triggers undefined behaviour, the server can do anything it likes _for that client_, but must not adversely affect other clients that did not trigger undefined behaviour.