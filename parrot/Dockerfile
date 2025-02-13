FROM scratch

COPY parrot /parrot
ENV PARROT_PORT=80
ENV PARROT_TRACE=true
EXPOSE 80

HEALTHCHECK --timeout=3s --start-interval=10ms --start-period=5s CMD [ "/parrot", "health" ]

ENTRYPOINT [ "/parrot" ]
