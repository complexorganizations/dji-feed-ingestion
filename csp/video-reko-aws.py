import boto3  # Import the boto3 library to access AWS services
import time  # Import the time library for waiting purposes


# Function to determine if two cars are close in time and space
def is_car_close(car1, car2, time_threshold, spatial_threshold):
    # Calculate the time difference between two cars
    time_diff = abs(car1['Timestamp'] - car2['Timestamp'])
    if time_diff > time_threshold:  # Check if the time difference is greater than the time threshold
        return False
    box1 = car1['BoundingBox']  # Get the bounding box for car1
    box2 = car2['BoundingBox']  # Get the bounding box for car2
    # Calculate the differences in the x and y coordinates of the centers of the two bounding boxes
    x_diff = abs((box1['Left'] + box1['Width'] / 2) -
                 (box2['Left'] + box2['Width'] / 2))
    y_diff = abs((box1['Top'] + box1['Height'] / 2) -
                 (box2['Top'] + box2['Height'] / 2))
    # Check if the spatial differences are less than or equal to the spatial threshold
    return x_diff <= spatial_threshold and y_diff <= spatial_threshold


# Function to count the cars in a video using Amazon Rekognition
def count_cars(video_file, bucket_name):
    # Create an Amazon Rekognition client
    rekognition = boto3.client('rekognition')
    response = rekognition.start_label_detection(  # Start the label detection job
        Video={
            'S3Object': {
                'Bucket': bucket_name,  # Specify the S3 bucket name
                'Name': video_file  # Specify the video file name
            }
        },
        MinConfidence=80  # Set the minimum confidence level for detected labels
    )
    job_id = response['JobId']  # Get the job ID from the response
    # Print the job ID
    print(f'Started label detection job with JobId: {job_id}')
    # Wait for the job to complete or fail
    while True:
        response = rekognition.get_label_detection(
            JobId=job_id)  # Get the current status of the job
        status = response['JobStatus']  # Get the job status from the response
        if status == 'SUCCEEDED':  # If the job is successful, break the loop
            break
        # If the job fails or times out, print the error and return None
        elif status == 'FAILED' or status == 'TIMED_OUT':
            print(f'Job failed with status: {status}')
            return None
        else:  # If the job is still in progress, wait for a while and then check again
            print(f'Waiting for job to complete... (Current status: {status})')
            time.sleep(10)
    labels = response['Labels']  # Get the labels from the response
    detected_cars = []  # Initialize an empty list to store detected cars
    time_threshold = 1000  # Set the time threshold in milliseconds
    # Set the spatial threshold (a value between 0 and 1)
    spatial_threshold = 0.05
    # Iterate through the labels to find cars
    for label in labels:
        if label['Label']['Name'] == 'Car':  # If the label is a car
            # Iterate through the instances of the car label
            for instance in label['Label']['Instances']:
                car = {
                    # Store the bounding box of the car instance
                    'BoundingBox': instance['BoundingBox'],
                    # Store the timestamp of the car instance
                    'Timestamp': label['Timestamp']
                }
                is_new_car = True  # Assume the detected car is a new car
                # Iterate through previously detected cars to check if the current car is already detected
                for detected_car in detected_cars:
                    if is_car_close(car, detected_car, time_threshold, spatial_threshold):
                        is_new_car = False  # If the car is close to a detected car, set is_new_car to False
                        break
                if is_new_car:  # If the detected car is indeed a new car
                    # Add the new car to the detected_cars list
                    detected_cars.append(car)
    return len(detected_cars)  # Return the number of detected cars


# Main function to execute the script
if __name__ == '__main__':
    video_file = 'output.mp4'  # Path to the video file
    bucket_name = 'test-video-feed-bucket'  # S3 bucket name
    # Call the count_cars function to count the cars in the video
    car_count = count_cars(video_file, bucket_name)
    # Print the number of cars detected in the video
    print(f'Number of cars detected: {car_count}')
