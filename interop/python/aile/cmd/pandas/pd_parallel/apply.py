from multiprocessing import Pool, cpu_count
import pandas as pd
import numpy as np

class Applier(object):
    def __init__(self, function, axis=1):
        self.f = function
        self.axis = axis
    def __call__(self, df):
        return df.apply(self.f, self.axis)
        
def applyParallel(df, func, cores=cpu_count()):
    split_df = np.array_split(df, cores)
    with Pool(cores) as pool:
        ret_list = pool.map(Applier(func), split_df)
    return pd.concat(ret_list)