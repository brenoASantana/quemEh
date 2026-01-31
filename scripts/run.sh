# Entra no front, builda
cd frontend
npm run build

# Apaga a dist antiga do back e copia a nova
rm -rf ../backend/dist
cp -r dist ../backend/

# Roda o back
cd ../backend
go run .