# HOOK :leftwards_arrow_with_hook:


## Examples

```sh
hook -args -cmd 'bundle exec rspec'
```

will run `bundle exec rspec <path to file which changed>`


```sh
hook -cmd 'bundle exec rspec' -dir ~/proj/yoloapp
```

will run `bundle exec rspec` for any file change in `~/proj/yoloapp`


## Bugs

- it doesn't fully work yet for some commands. (environment problems etc etc)

## License

MIT
