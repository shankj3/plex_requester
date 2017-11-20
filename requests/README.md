docker url 
```
sudo docker run --restart=always -p 80:8080 -d -v /home/ec2-user/perplexed_data:/tmp/perplexed_data --name perplexed jessishank/perplexed -movies_file=/tmp/perplexed_data/movies.json -finished_file=/tmp/perplexed_data/finished.json
```