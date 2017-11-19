-- jobs ratio
SELECT count(a.userId) as total, a.answer FROM answers a
  WHERE a.questionId = "job"
  GROUP BY a.answer

-- location ratio
SELECT count(a.userId) as total, a.answer FROM answers a
  WHERE a.questionId = "island"
  GROUP BY a.answer


