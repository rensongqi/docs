
# 通过shell脚本实现弱密码检查

1. 对比/etc/shadow文件中guest用户是否存在弱密码，支持主流的md5、blowfish、sha-256、sha-512、yescrypt算法
2. 存在弱密码则修改密码，并把新密码加盐后使用公钥加密存放在/etc/product/.passwd文件中
3. 如果密码修改失败，则用户密码回退至旧密码
4. 每隔15天密码将自动重置以防止密码泄露

```bash
#!/bin/bash

# weak password example
dict=("!@#"
"!@#$%^"
"!@#$%^&*"
"!@#$%^&*()_+"
"%resu%"
"%tsoh%"
"%user%"
"%user%!"
"%user%!@#"
"%user%."
"%user%00"
"%user%111"
"%user%12"
"%user%123"
"%user%1234"
"%user%888"
"%user%@"
"%user%ab"
"%user%abc"
"%user%abcd"
"(%user%)"
".%user%."
"000000"
"00000000"
"111111"
"11111111"
"111222"
"123"
"123!@#"
"123123"
"123321"
"1234"
"12345"
"123456"
"1234567"
"12345678"
"123456789"
"1234567890-="
"123qwe"
"222222"
"333333"
"444444"
"555555"
"654321"
"6543210"
"6666"
"666666"
"66666666"
"666888"
"816357"
"8888"
"888888"
"88888888"
"999999"
"99999999"
"Admin123"
"Guest123"
"^%user%^"
"1234567890-="
"aaaaaa"
"ab"
"abc"
"abc123"
"abcd"
"abcd1234"
"abcde"
"abcdef"
"abcdefg"
"abcdefgh"
"abcdefghi"
"admin"
"asdfghjkl"
"asdfjkl;"
"exit"
"guest"
"iloveyou"
"master"
"password"
"qazwsx"
"qq123456"
"qweasd"
"root"
"taobao"
"test"
"test123"
"wang1234"
"woaini"
"xiaoming"
"zzzzzz"
"xxxxxxxx"
" ")

pub_key="-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6F+9rKUlGqEo02NOkRGK\nlsgpMVS5mQKmLCpwOvQzGjOUZM08aC1q/Ygt9iXS/FdSrcB7b3cl99sMizZCdHRz\nHEWH9mmeArjY60Nm4QL/cFD8YTm8z0SAcwJZQT1cPGG5YkoQeTnwYLoxjVx75fzt\nsv2iTscLQZ+vXu1oO1k7Bh0bK8FLEsifp9o5a5CqpvYJTJiCdX0xAHfFu8HBsVlH\nklow+Eghc0PPa16pJg0o8VNcHMZdpsHL0k134Mm0rehbuYpuAL2s8l69BiU/44+Y\nfK+hm8frYjVoEBOt2EBKeFbeeHM0YotKkogKnFQ87vTrrvBcPW7YDiKO+5FVUqQY\nQQIDAQAB\n-----END PUBLIC KEY-----"
pub_path="/tmp/client.pub"
salt_key="pAOx20AkCxCrvZqv"
pass_dir=/etc/product
pass_file=".passwd"
pass_file_path="$pass_dir/$pass_file"

notify_iot() {
    echo "notifying iot ..."
}

change_pass() {
    user=$1
    old_password=$2
    new_password=$3

    if [[ ! -d $pass_dir ]]; then
        mkdir -p $pass_dir
    fi
    # change passwd
    echo "$user:$new_password" | chpasswd

    # check whether the password is successfully changed
    if [ $? -eq 0 ]; then
        echo "$new_password$salt_key" > "$pass_file_path-temp"
        encrypt_pass_file
        echo "Password for $1 has been successfully changed."
    else
        # if the change fails, roll back the old password
        echo "$user:$old_password" | chpasswd
    fi
}

change_current_user_pass() {
    # get file modify time
    modification_time=$(stat -c %Y "$pass_file_path")

    # get current time
    current_time=$(date +%s)

    # calculate the time difference
    time_difference=$((current_time - modification_time))

    # check if the password file has been modified for more than 15 days
    if [ $time_difference -gt $((15 * 24 * 3600)) ]; then
        old_pass_en=$(<$pass_file_path-temp)
        old_pass=$(echo $old_pass_en | cut -c 1-12)
        generate_rand_pass
        change_pass "guest" $old_pass $new_pass
    fi
}

generate_rand_pass() {
    pass_length=12
    pass=$(openssl rand -base64 48 | tr -dc 'a-zA-Z0-9')
    new_pass=${pass:0:$pass_length}
}

check_pass() {
    host=$(hostname -s)
    tsoh=$(echo $host | rev)

    while IFS=':' read -r username password _; do
        # sha256 sha-512 ...
        if [[ ! $password =~ [!!|*] ]] && [[ ! $password =~ ^\$y\$j9T ]]; then
            IFS='$' read -r _ const salt hash <<< "$password"
            for d in "${dict[@]}"; do
                t="$d"
                resu=$(echo "$username" | rev)

                t="${t//\%user\%/$username}"
                t="${t//\%resu\%/$resu}"

                t="${t//\%host\%/$host}"
                t="${t//\%tsoh\%/$tsoh}"

                hashed=$(perl -e "print crypt('$t', '\$$const\$$salt\$')")

                if [ "$hashed" == "$password" ] && [ "$username" == "guest" ]; then
                    generate_rand_pass
                    change_pass $username $t $new_pass
                    echo "$username [$t]"
                    break
                fi
            done
        # yescrypt algorithm
        elif [[ $password =~ ^\$y\$j9T ]]; then
            IFS='$' read -r _ const str salt hash <<< "$password"
            for d in "${dict[@]}"; do
                t="$d"
                resu=$(echo "$username" | rev)

                t="${t//\%user\%/$username}"
                t="${t//\%resu\%/$resu}"

                t="${t//\%host\%/$host}"
                t="${t//\%tsoh\%/$tsoh}"

                hashed=$(perl -e "print crypt('$t', '\$$const\$$str\$$salt\$')")

                if [ "$hashed" == "$password" ] && [ "$username" == "guest" ]; then
                    generate_rand_pass
                    change_pass $username $t $new_pass
                    echo "$username [$t]"
                    break
                fi
            done
        fi
    done < /etc/shadow
}

encrypt_pass_file() {
    # save pub_key to file
    echo -e $pub_key > $pub_path

    # use public key to encrypt
    openssl rsautl -encrypt -pubin -inkey $pub_path -in $pass_file_path-temp -out $pass_file_path

    # use private key to decrypt
    # openssl rsautl -decrypt -inkey client.key -in $pass_file_path -out $pass_file_path-decode

    # delete pub_key file
    rm -f $pub_path
}

func() {
    echo "Usage: check.sh [-a|-c|-h]"
    echo ""
    echo "    -a, automatically check for weak passwords and change them"
    echo "    -c, changing the password of an existing user"
    exit -1
}
 
while getopts ':ach' OPT; do
    case $OPT in
        a) check_pass;;
        c) change_current_user_pass;;
        h) func;;
        ?) func;;
    esac
done
```