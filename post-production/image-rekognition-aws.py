import boto3


# Get the bucket name from a s3 URI
def get_bucket_name_from_s3_uri(s3_uri):
    return s3_uri.split("/")[2]


# Get the path from a s3 URI
def get_path_from_s3_uri(s3_uri):
    return "/".join(s3_uri.split("/")[3:])


# Get the response from Rekognition
def get_response(s3URI):
    rekognition = boto3.client("rekognition")
    response = rekognition.detect_labels(
        Image={
            "S3Object": {
                "Bucket": get_bucket_name_from_s3_uri(s3URI),
                "Name": get_path_from_s3_uri(s3URI),
            }
        }
    )
    return response


# Get all the file paths of all the objects in a bucket
def get_all_file_path_in_s3_bucket(s3URI):
    s3 = boto3.client("s3")
    response = s3.list_objects_v2(Bucket=get_bucket_name_from_s3_uri(
        s3URI), Prefix=get_path_from_s3_uri(s3URI))
    # Return the file paths of all the objects in the bucket
    return [object["Key"] for object in response["Contents"]]


def main():
    # Get the response from Rekognition
    response = get_response(
        "s3://dji-live-stream-feed-data-0/media/dji-stream-0//113669145637_dji-stream-0_1675350168417_db41d0f6-ccf1-4665-aa8d-de74561aaa64.jpg")

    # Print the response
    print(response)

    # Print the confidence and name of each label
    for label in response["Labels"]:
        print(label["Name"], label["Confidence"])

    # Get all the file paths of all the objects in a bucket
    objects = get_all_file_path_in_s3_bucket(
        "s3://dji-live-stream-feed-data-0/")
    print(objects)


main()
