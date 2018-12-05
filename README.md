# godd - Go OvS Datapath Demystifier
`godd` is a CLI tool that makes the OvS' output for the datapath flows human readable
![img1](https://gitlab.com/rdodin/pics/wikis/uploads/63ac9c98c06fd16bddf2d60c3711b3f3/image.png)

It accepts flows information in the standard input and parses it out to a YAML structured output.

## Usage scenarios
`godd`'s primary goal is to make Nuage openvswitch datapath troubleshooting, verification and monitoring easier to human beings.  
Despite that I created `godd` specifically for Nuage flavored OvS, it should work well on a vanilla OvS instances.

### 1 Running on an OvS-powered device
You can download and run `godd` on the OvS device itself:
```
ovs-dpctl dump-flows -m | ~/godd
```

### 2 Running outside of OvS-powered device
Its not always convenient to download the tool to each and every OvS switch, sometimes you already have the flows output that you can transform with `godd`.

Since the tool expects the data on the standard input you can either `cat` a file to it:
```
cat myFlowsData.txt | ~/godd
```

or simply `echo` the flows you have:
```
echo 'recirc_id(0),tunnel(tun_id=0xc882d,src=10.10.99.99,dst=10.10.1.11,ttl=63,flags(-df-csum+key)),in_port(10),skb_mark(0x20000000),eth(src=5a:eb:de:fe:6b:44,dst=00:00:0a:0a:01:0b),eth_type(0x0800),ipv4(src=172.254.0.59,dst=172.254.0.12,proto=17,ttl=63,frag=no),udp(src=19995,dst=19995), packets:4, bytes:360, used:2.418s, actions:set(eth(src=68:54:ed:00:ec:9e,dst=68:54:ed:00:f1:0e)),set(ipv4(dst=172.254.0.12,ttl=62)),14' | ~/godd
```

## Download
The tool is distributed as a single binary for Linux/Mac/Windows OSes.

Check the links for download the latest (`v.0.0.1`) version for each supported OS. Feel free to correct the download path to someplace in your `$PATH` if you wish so.

**Linux**:
```
curl -kL https://github.com/hellt/godd/releases/download/v0.0.1/godd_linux > ~/godd && chmod a+x ~/godd
```

**Mac**:
```
curl -kL https://github.com/hellt/godd/releases/download/v0.0.1/godd_darwin > ~/godd && chmod a+x ~/godd
```

**Windows**:
```
curl -kL https://github.com/hellt/godd/releases/download/v0.0.1/godd_windows > ~/godd && chmod a+x ~/godd
```

## What are the godd benefits?
Despite being nice and output the flow info in a nicely looking YAML (opinionated) `godd` sports a few additional enhancements:

1. easy flow distinction: each flow will be separated with two newlines and a header in green color with the flow seq. number. Since the YAML output is verbose, the ease of flow distinction makes it a breeze to scroll over the flows and find the right one
2. grouping related data: the flow data that nests in a common place will be grouped together. I.e. all the `set()` actions will be grouped under a single `set:` parent, making the result easier to parse with a pair of eyes.
3. tunnel_id conversion: in Nuage, hex encoded tunnnel ID has a decimal counterpart that is present on the controllers. `godd` augments hex tun_id field with the decimal representation.

## What next?
There are a lot of improvements that could be made. Some of them are Nuage related, some are of a general kind:

1. output in a JSON/table/condensed formats
2. receive openflow ports data to unwrap the port names
3. add a hash for the 5-tuple data to make a correlation between the flows on a different switches
4. and more

## Contribution
The easiest way to contribute is to file an issue with the dump-flow data that `godd` failed to parse, this will give me more raw data to analyze and make `godd` more generic.

## Author
Roman Dodin - [@ntdvps](https://twitter.com/ntdvps)