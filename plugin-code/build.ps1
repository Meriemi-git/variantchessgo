docker run --rm -w "/builder" -v "${PWD}:/builder" heroiclabs/nakama-pluginbuilder:2.11.1 build -buildmode=plugin -trimpath -o ./modules/plugin_code.so
Copy-Item -path ./modules/plugin_code.so -destination C:\Users\Doremus\docker\nakama\modules
cd C:\Users\Doremus\docker\nakama
docker-compose restart