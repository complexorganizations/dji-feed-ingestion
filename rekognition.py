import boto3

def init():
    s3 = boto3.resource("s3")
    rekognition = boto3.client("rekognition")


init()

def main():
    # Get the response from Rekognition
    response = get_response("bucket", "key")

main()

# Get the response from Rekognition
def get_response(bucket, key):
    response = rekognition.detect_labels(
        Image={
            "S3Object": {
                "Bucket": bucket,
                "Name": key
            }
        }
    )
    return response

