epigrams.txt:
	curl -s https://www.cs.yale.edu/homes/perlis-alan/quotes.html | pup p | html2text -width 256 > $@

clean:
	rm -f epigrams.txt
