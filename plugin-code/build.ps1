docker run --rm -w "/builder" -v "${PWD}:/builder" heroiclabs/nakama-pluginbuilder build -buildmode=plugin -trimpath -o ./modules/plugin_code.so
Copy-Item -path ./modules/plugin_code.so -destination B:\Dev\docker\nakama\modules
cd B:\Dev\docker\nakama
docker-compose restart