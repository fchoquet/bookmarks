FROM centurylink/ca-certs
EXPOSE 8080
ADD docs /docs
ADD templates /templates
ADD bin/bookmarks /
CMD ["/bookmarks"]
