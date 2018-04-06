import csv
import re
import pandas as pd


def reg_test(name):

    reg_result = ''

    with open(name, 'r') as csvfile:
        reader = csv.reader(csvfile)
        for row in reader:
            row = re.sub('[^A-Za-z0-9]+', '', str(row))
            reg_result += row + ','

    return reg_result

def cleaner(fileName, colName): 
    data = pd.read_csv(fileName)   
    with open("stopwords.txt", 'r') as f:
        stopwords = f.read().split("\n")
    stopwords = [re.sub(r'[^\w\s\d]','',s.lower()) for s in stopwords]

    vec = data[colName]
    # Lowercase, then replace any non-letter, space, or digit character
    new_vec = [re.sub(r'[^\w\s\d]','', h.lower()) for h in vec]
    # Replace sequences of whitespace with a space character.
    new_vec = [re.sub("\s+", " ", h) for h in new_vec]
    # new_vec = [(w for w in row if w not in stopwords) for row in new_vec]
    data[colName] = new_vec
    
    print data[colName][0]
    data.to_csv( fileName + '_clean.csv')

def combiner(outFile):
    final = ''
    jobs = pd.read_csv('jobs.csv')
    cols = jobs.columns.values.tolist()
    # print jobs.dtypes
    for col in jobs.columns.values.tolist():
        vec = jobs[col]
        # # Lowercase, then replace any non-letter, space, or digit character
        try:
            new_vec = [re.sub(r'[^\w\s\d]','', h.lower()) for h in vec]
        except:
            new_vec = [re.sub(r'[^\w\s\d]','', str(h).lower()) for h in vec]
            continue
        # # Replace sequences of whitespace with a space character.
        new_vec = [re.sub("\s+", " ", h) for h in new_vec]
        jobs[col] = new_vec
    
    for i, r in jobs.iterrows():
        # print r
        for col in jobs.columns.values.tolist():
            try:
                final = final + r[col]
            except:
                final = final + str(r[col])
    transcripts = pd.read_csv('transcripts.csv')
    vec = transcripts['transcript']
    new_vec = [re.sub(r'[^\w\s\d]','', h.lower()) for h in vec]
    # # Replace sequences of whitespace with a space character.
    new_vec = [re.sub("\s+", " ", h) for h in new_vec]
    transcripts['transcript'] = new_vec
    for item in transcripts['transcript']:
        final = final + item
    print final[:50] + final[50000:]
    f = open(outFile,"w+")
    f.write("%s" % (final))
combiner('combined.csv')

