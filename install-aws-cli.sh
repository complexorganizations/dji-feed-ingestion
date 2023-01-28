apt-get install curl unzip -y
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
aws configure set aws_access_key_id ${{ secrets.AWS_ACCESS_KEY }}
aws configure set aws_secret_access_key ${{ secrets.AWS_SECRET_ACCESS_KEY }}
