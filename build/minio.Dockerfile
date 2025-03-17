# Используем официальный образ MinIO
FROM minio/minio:latest

# (Опционально) Добавьте свои настройки или конфигурации, если нужно
# Например, копирование конфигурационных файлов или скриптов
# COPY ./config.json /root/.minio/config.json

# Указываем точку входа
ENTRYPOINT ["minio"]
CMD ["server", "/data", "--console-address", ":9001"]