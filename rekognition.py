import boto3


# Get the bucket name from a s3 URI
def get_bucket_name_from_s3_uri(s3_uri):
    return s3_uri.split("/")[2]


# Get the path from a s3 URI
def get_path_from_s3_uri(s3_uri):
    return "/".join(s3_uri.split("/")[3:])


# Get the response from Rekognition
def get_response(s3URI):
    rekognition = boto3.client("rekognition", region_name="us-east-1")
    response = rekognition.detect_labels(
        Image={
            "S3Object": {
                "Bucket": get_bucket_name_from_s3_uri(s3URI),
                "Name": get_path_from_s3_uri(s3URI),
            }
        }
    )
    return response


def main():
    # Get the response from Rekognition
    response = get_response(
        "s3://dji-live-stream-feed-data-0/media/dji-stream-0//113669145637_dji-stream-0_1675350588904_b574d515-c67b-49af-8c74-48fd302bd3c0.jpg")
    print(response)


main()
