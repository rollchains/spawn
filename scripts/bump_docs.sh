# bumps docs versions for spawn

OLD_VERSION=v0.50.14
NEW_VERSION=v0.50.15

findAndReplace() {
    find . -type f -name "$1" -not -path "*node_modules*" -exec sed -i "$2" {} \;
}

findAndReplace "*.md" "s/$OLD_VERSION/$NEW_VERSION/g"
findAndReplace "install.sh" "s/$OLD_VERSION/$NEW_VERSION/g"
