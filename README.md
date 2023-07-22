# totp

Simple cli tool for getting TOTP codes.

Add a couple of new secrets:

```
totp add facebook RW***MN*****U7XCZIVZ
totp add work OUVCQ******YAOM
```

Now when you need to get the code just run

```
totp get work # now you can paste the code
```

You can set a default secret in the config `~/Library/Application Support/totp.yaml`
```
default: work
secrets:
    facebook: RW***MN*****U7XCZIVZ
    work: OUVCQ******YAOM
```

That will give you the code for `work` by running `totp` without arguments


# How to transfer you accounts from Google Authenticator

- open the Google Authenticator app
- choose `Transfer accounts` in the menu
- then Export
- once you have a QR code - find a way to save it as an image  
  use you laptop's camera for example

- brew install zbar
- decode QR code with `zbarimg qr1.jpg > qr1.txt`
- run `./totp decode < qr1.txt`
- you will get a list of secrets
- now you can transfer those accounts into any application

# Limitations

Works only on mac