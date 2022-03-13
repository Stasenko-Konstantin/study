module Main where

import Control.Parallel
import Data.List

pfold :: (Num a, Enum a) => (a -> a -> a) -> a -> [a] -> a
pfold f z []  = z
pfold f z [a] = a
pfold f z xs  = (ys `par` zs) `pseq` (ys `f` zs) where
  len = length xs
  (ys', zs') = splitAt (len `div` 2) xs
  ys = pfold f z ys'
  zs = pfold f z zs'

-- foldl :: (b -> a -> b) -> b -> t a -> b
-- foldl f z []     = z
-- foldl f z (x:xs) = foldl f (f z x) xs

f x = x

a = 0.0

b = 1.0

n = 100000000

ns = [1 .. n] :: [Double]

h = (b - a) / n

main :: IO ()
main = do
  print $
    show $
      0.5 * (f a + f b) * h
        + pfold
          (\r i -> r + f (a + i * h) * h)
          0.0
          ns
