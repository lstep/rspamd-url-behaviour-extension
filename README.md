# URL Behaviour extension for RSPAMd

Goal of this project is to implement an extension for RSPAMd that is capable of examining all the urls in an email with the following features:

- [ ] ML analysis of the presented url (ultra long url, domains, subdomains, known fake name one lowest subdomain (paypal, etc))
- [ ] Analysis of the domain + SSL
    - known in blacklist/reputation ? (original and final redirected)
        - Use 1 million Alexa top ranking (dead?)
        - Use Cisco Umbrella (https://s3-us-west-1.amazonaws.com/umbrella-static/index.html)
        - Majestic million (https://downloads.majesticseo.com/majestic_million.csv)
    - when registered + SSL cert creation?
- [ ] Check for redirections, and if too many (> 2 for example), will trigger 
- [ ] send url to a sandbox and examine its content (AI, visual check)
- [ ] Analysis of the final page:
    - number of broken links
    - AI analysis of visual representation
    - Check the Domain Age and Ownership (whois search)
    - Check SSL certificate origin (and date creation)


# References
- Jeffry Sleddens's rspamd plugin: https://github.com/jeffrysleddens/rspamd-bitcoinabuse-plugin
- SwissCenter made or modded plugins for rspamd: https://github.com/sriccio/rspamd-plugins
- Phishing URL Detection with ML: https://towardsdatascience.com/phishing-domain-detection-with-ml-5be9c99293e5
