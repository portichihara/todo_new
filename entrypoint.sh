echo "Waiting for postgres..."

while ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER
do
  echo "Waiting for postgres to be ready..."
  sleep 2
done

echo "PostgreSQL is ready!"

echo "Current directory: $(pwd)"
echo "Directory contents:"
ls -la
echo "Templates directory contents:"
ls -la templates/

./app