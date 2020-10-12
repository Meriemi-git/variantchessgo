
docker run --rm -w "/builder" -v "${PWD}:/builder" heroiclabs/nakama-pluginbuilder build -buildmode=plugin -trimpath -o ./plugin-code/modules/plugin_code.so
Copy-Item -path ./plugin-code/modules/plugin_code.so -destination B:\Dev\docker\nakama\modules
Set-Location B:\Dev\docker\nakama
docker-compose restart