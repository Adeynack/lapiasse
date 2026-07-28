User.create_with(display_name: "Joe", password: "joe").find_or_create_by!(email: "joe@example.com")
User.create_with(display_name: "Proud Mary", password: "mary").find_or_create_by!(email: "mary@example.com")
User.create_with(display_name: "Vlad the Impaler", password: "vlad").find_or_create_by!(email: "vlad@example.com")
