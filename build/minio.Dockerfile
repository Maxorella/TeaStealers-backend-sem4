FROM minio/minio:latest
ENTRYPOINT ["minio"]
CMD ["server", "/data", "--console-address", ":9001"]