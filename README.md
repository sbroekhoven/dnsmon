# dnsmon

This `dnsmon` tool get's information about a domainname from DNS, and looks for changes with the previous lookup.

To start with, I only want to look for changes in:

* The serial in the SOA record to see if there was any update
* MX records for mailserver changes
* NS records to see if the nameserver has changed

