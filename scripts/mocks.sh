#!/usr/bin/env bash

echo "Generate mocks..."
mockery --all --inpackage --with-expecter=true \
  --case=underscore --inpackage-suffix --testonly \
  --dir=internal/domain
