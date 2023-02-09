import boto3


# Get the bucket name from a s3 URI
def get_bucket_name_from_s3_uri(s3_uri):
    return s3_uri.split("/")[2]


# Get the response from Rekognition
def get_response(s3URI):
    rekognition = boto3.client("rekognition", region_name="us-east-1")
    response = rekognition.detect_labels(
        Image={
            "S3Object": {
                "Bucket": get_bucket_name_from_s3_uri(s3URI),
                "Name": s3URI
            }
        }
    )
    return response


def main():
    # Get the response from Rekognition
    response = get_response("s3://dji-live-stream-feed-data-0/media/dji-stream-0//113669145637_dji-stream-0_1675351204694_f8953f18-b68f-469e-bd16-7b972f549c62.jpg")
    print(response)


main()
