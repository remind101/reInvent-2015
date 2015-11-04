FROM nginx

COPY nginx.conf /etc/nginx/nginx.conf
COPY index.html acme.jpg /usr/share/nginx/html/
ENV PORT=8080

CMD ["nginx"]
