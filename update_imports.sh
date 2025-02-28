#!/bin/bash

# Update SQLBoiler imports
find . -type f -name "*.go" -not -path "./models/*" | xargs sed -i '' 's|github.com/volatiletech/sqlboiler/boil|github.com/volatiletech/sqlboiler/v4/boil|g'
find . -type f -name "*.go" -not -path "./models/*" | xargs sed -i '' 's|github.com/volatiletech/sqlboiler/queries|github.com/volatiletech/sqlboiler/v4/queries|g'
find . -type f -name "*.go" -not -path "./models/*" | xargs sed -i '' 's|github.com/volatiletech/sqlboiler/drivers|github.com/volatiletech/sqlboiler/v4/drivers|g'
find . -type f -name "*.go" -not -path "./models/*" | xargs sed -i '' 's|github.com/volatiletech/sqlboiler/strmangle|github.com/volatiletech/sqlboiler/v4/strmangle|g'

# Update null imports
find . -type f -name "*.go" -not -path "./models/*" | xargs sed -i '' 's|github.com/volatiletech/null|github.com/volatiletech/null/v8|g'

echo "Imports updated successfully!" 