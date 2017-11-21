FROM acoshift/go-scratch

ADD postlist /
ADD template /template
ADD table.sql /

EXPOSE 8080
ENTRYPOINT ["/postlist"]

