# Go tests
guard :shell do
  watch /.*\.go/ do |m|
    if m[0].end_with? '_test.go'
        n m[0], "Test has changed."
    else
        n m[0], "has changed."
    end
    file = File.dirname(m[0])
    `cd #{file} && go test`
  end
end
