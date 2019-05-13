# Refresh "model"
rm -rf ./vendor/github.com/ekara-platform/model/*.go
mkdir -p ./vendor/github.com/ekara-platform/model/
cp ../model/*.go  ./vendor/github.com/ekara-platform/model/


# Refresh "engine"
rm -rf ./vendor/github.com/ekara-platform/engine/*.go
cp ../engine/*.go  ./vendor/github.com/ekara-platform/engine/

rm -rf ./vendor/github.com/ekara-platform/engine/ansible/*.go
mkdir -p ./vendor/github.com/ekara-platform/engine/ansible/
cp ../engine/ansible/*.go  ./vendor/github.com/ekara-platform/engine/ansible/

rm -rf ./vendor/github.com/ekara-platform/engine/component/*.go
mkdir -p ./vendor/github.com/ekara-platform/engine/component/
cp ../engine/component/*.go  ./vendor/github.com/ekara-platform/engine/component/

rm -rf ./vendor/github.com/ekara-platform/engine/ssh/*.go
mkdir -p ./vendor/github.com/ekara-platform/engine/ssh/
cp ../engine/ssh/*.go  ./vendor/github.com/ekara-platform/engine/ssh/

rm -rf ./vendor/github.com/ekara-platform/engine/util/*.go
mkdir -p ./vendor/github.com/ekara-platform/engine/util/
cp ../engine/util/*.go  ./vendor/github.com/ekara-platform/engine/util/

rm -rf ./vendor/github.com/ekara-platform/engine/component/scm/*.go
mkdir -p ./vendor/github.com/ekara-platform/engine/component/scm/
cp ../engine/component/scm/*.go  ./vendor/github.com/ekara-platform/engine/component/scm/

rm -rf ./vendor/github.com/ekara-platform/engine/component/scm/file/*.go
mkdir -p ./vendor/github.com/ekara-platform/engine/component/scm/file/
cp ../engine/component/scm/file/*.go  ./vendor/github.com/ekara-platform/engine/component/scm/file/

rm -rf ./vendor/github.com/ekara-platform/engine/component/scm/git/*.go
mkdir -p ./vendor/github.com/ekara-platform/engine/component/scm/git/
cp ../engine/component/scm/git/*.go  ./vendor/github.com/ekara-platform/engine/component/scm/git/

