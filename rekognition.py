import boto3


# Get the response from Rekognition
def get_response(bucket, key):
    rekognition = boto3.client("rekognition", region_name="us-west-1")
    response = rekognition.detect_labels(
        Image={
            "S3Object": {
                "Bucket": bucket,
                "Name": key
            }
        }
    )
    return response


def main():
    # Get the response from Rekognition
    response = get_response("bucket", "key")
    print(response)


main()
