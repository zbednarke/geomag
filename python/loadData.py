import json
import numpy as np
import tensorflow as tf
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler

# Optionally, create a Python class mirroring the Go struct
class Dataset:
    def __init__(self, latitude, longitude, altitude, bx, by, bz):
        self.latitude = latitude
        self.longitude = longitude
        self.altitude = altitude
        self.bx = bx
        self.by = by
        self.bz = bz

def load_dataset_from_json(file_path):
    with open(file_path, 'r') as f:
        data = json.load(f)
    # `data` is a dict with keys: "latitude", "longitude", etc.
    ds = Dataset(**data)
    return ds

def build_model():
    """
    Build a small feed-forward neural network in TensorFlow/Keras 
    to predict longitude from magnetic field components (Bx, By, Bz).
    """
    model = tf.keras.Sequential([
        tf.keras.layers.Dense(64, activation='relu', input_shape=(3,)),
        tf.keras.layers.Dense(64, activation='relu'),
        tf.keras.layers.Dense(1)  # Single output: predicted longitude
    ])
    model.compile(
        optimizer='adam',
        loss='mean_squared_error',
        metrics=['mean_absolute_error']
    )
    return model

def main():

    gpus = tf.config.experimental.list_physical_devices('GPU')
    if gpus:
        print("GPUs detected:")
        for gpu in gpus:
            print(f"- {gpu}")
    else:
        print("No GPUs detected. TensorFlow is using the CPU.")


    # 1. Load data
    ds = load_dataset_from_json("./training/equatorDataset.json")

    # Convert to NumPy arrays (assuming they're lists in the JSON)
    bx = np.array(ds.bx, dtype=np.float32)
    by = np.array(ds.by, dtype=np.float32)
    bz = np.array(ds.bz, dtype=np.float32)
    longitudes = 180 / np.pi * np.array(ds.longitude, dtype=np.float32)
    
    # 2. Prepare features (X) and target (y)
    X = np.stack([bx, by, bz], axis=1)  # shape: (num_samples, 3)
    y = longitudes.reshape(-1, 1)       # shape: (num_samples, 1)

    idxs = np.random.permutation(len(X))
    X = X[idxs]
    y = y[idxs]

    # 3. Split data into train/test sets
    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42
    )

    # 4. Scale inputs and outputs using StandardScaler
    X_scaler = StandardScaler()
    y_scaler = StandardScaler()

    # Fit on the training data, transform both train and test
    X_train_scaled = X_scaler.fit_transform(X_train)
    X_test_scaled = X_scaler.transform(X_test)

    y_train_scaled = y_scaler.fit_transform(y_train)
    y_test_scaled = y_scaler.transform(y_test)

    # 5. Build and train the model on scaled data
    model = build_model()
    print(model.summary())

    history = model.fit(
        X_train_scaled, 
        y_train_scaled,
        epochs=50,            # Adjust as desired
        batch_size=32,        # Adjust as desired
        validation_split=0.2, # Uses a portion of the training set for validation
        verbose=1
    )

    # 6. Evaluate on the test set (scaled)
    test_loss, test_mae = model.evaluate(X_test_scaled, y_test_scaled, verbose=0)
    print(f"\nTest MSE (loss): {test_loss:.4f}, Test MAE (scaled): {test_mae:.4f}")

    # 7. Make predictions on the test set in scaled space, then invert
    preds_scaled = model.predict(X_test_scaled)
    preds_unscaled = y_scaler.inverse_transform(preds_scaled)  # revert to original longitude space

    # If you want to see the "true" values in the original scale, use y_scaler as well:
    y_test_unscaled = y_scaler.inverse_transform(y_test_scaled)

    # Compare a few predictions with the ground truth
    print("\nSample predictions vs. actual (longitude):")
    for i in range(10):
        print(f"Predicted: {preds_unscaled[i,0]:.5f} | Actual: {y_test_unscaled[i,0]:.5f} | diff: {np.pi / 180  * (preds_unscaled[i,0] - y_test_unscaled[i,0]):.5f}rads, or {6.378 * 1e3 * np.pi / 180 * (preds_unscaled[i,0] - y_test_unscaled[i,0]):.5f}km")
    
    all_errs_km = np.abs(6.378 * 1e3 * np.pi / 180 * (preds_unscaled[:,0] - y_test_unscaled[:,0]))
    print(f"\nMean abs error: {np.mean(all_errs_km):.5f} km")


    
if __name__ == "__main__":
    main()
